from reportlab.platypus import BaseDocTemplate, Frame, PageTemplate
from reportlab.platypus import Paragraph, Spacer, Image
from reportlab.platypus import SimpleDocTemplate
from reportlab.lib.units import mm
from reportlab.pdfbase.cidfonts import UnicodeCIDFont
from reportlab.pdfbase import pdfmetrics
from reportlab.lib.styles import ParagraphStyle, getSampleStyleSheet
from reportlab.lib import colors
from reportlab.platypus import Table, TableStyle
from datetime import datetime
import sys
import os

output_path = sys.argv[1]
text = sys.argv[2]
image_path = sys.argv[3]
staff_name = sys.argv[4]
logo_path = sys.argv[5]

# 日本語フォント登録
pdfmetrics.registerFont(UnicodeCIDFont('HeiseiMin-W3'))

# ドキュメント設定（余白固定）
doc = SimpleDocTemplate(
    output_path,
    rightMargin=20,
    leftMargin=20,
    topMargin=30,
    bottomMargin=30,
)

styles = getSampleStyleSheet()

normal_style = ParagraphStyle(
    'JapaneseNormal',
    parent=styles['Normal'],
    fontName='HeiseiMin-W3',
    fontSize=11,
    leading=15,
)

title_style = ParagraphStyle(
    'JapaneseTitle',
    parent=styles['Heading1'],
    fontName='HeiseiMin-W3',
    fontSize=16,
    leading=20,
    alignment=1,  # 中央揃え
)

elements = []

# =========================
# ヘッダー（ロゴ＋日付）
# =========================

today = datetime.now().strftime("%Y/%m/%d")

header_data = []

# ロゴ
if os.path.exists(logo_path):
    logo = Image(logo_path)
    logo.drawHeight = 15 * mm
    logo.drawWidth = 40 * mm
else:
    logo = Paragraph("", normal_style)

# 日付（右上）
date_paragraph = Paragraph(f"作成日：{today}", normal_style)

header_data.append([logo, date_paragraph])

header_table = Table(header_data, colWidths=[100 * mm, 60 * mm])
header_table.setStyle(TableStyle([
    ('ALIGN', (1, 0), (1, 0), 'RIGHT'),
    ('VALIGN', (0, 0), (-1, -1), 'MIDDLE'),
]))

elements.append(header_table)
elements.append(Spacer(1, 15 * mm))

# =========================
# タイトル
# =========================

elements.append(Paragraph("社内報告書", title_style))
elements.append(Spacer(1, 10 * mm))

# =========================
# 担当者
# =========================

elements.append(Paragraph(f"担当者：{staff_name}", normal_style))
elements.append(Spacer(1, 10 * mm))

# =========================
# 本文
# =========================

elements.append(Paragraph("【報告内容】", normal_style))
elements.append(Spacer(1, 5 * mm))

elements.append(Paragraph(text, normal_style))
elements.append(Spacer(1, 15 * mm))

# =========================
# 画像＋キャプション
# =========================

if os.path.exists(image_path):
    img = Image(image_path)
    img.drawWidth = 120 * mm
    img.drawHeight = 90 * mm
    elements.append(img)
    elements.append(Spacer(1, 5 * mm))
    elements.append(Paragraph("図1：施工予定箇所", normal_style))

# PDF生成
doc.build(elements)
