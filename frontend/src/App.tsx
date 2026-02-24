import React, { useState } from "react";
import { PhotoUpload } from "./components/PhotoUpload";
import type { Photo } from "./components/PhotoUpload";

import "./index.css";

export default function App() {
  const [entrancePhotos, setEntrancePhotos] = useState<Photo[]>([
    { file: null, preview: "" },
    { file: null, preview: "" },
  ]);
  const [hallwayPhotos, setHallwayPhotos] = useState<Photo[]>([
    { file: null, preview: "" },
    { file: null, preview: "" },
  ]);
  const [comment, setComment] = useState("");
  const [loading, setLoading] = useState(false);

  const handleGenerateReport = async () => {
    try {
      setLoading(true);
      const formData = new FormData();
      // entrancePhotos.forEach((p) => p.file && formData.append("text", "玄関に手すりをつけたい") && formData.append("image", p.file));
      // hallwayPhotos.forEach((p) => p.file && formData.append("text", "廊下に手すりをつけたい") && formData.append("image", p.file));

      entrancePhotos.forEach((p) => {
        if (p.file) {
          formData.append("text", "玄関に手すりをつけたい");
          formData.append("image", p.file);
        }
      });

      hallwayPhotos.forEach((p) => {
        if (p.file) {
          formData.append("text", "廊下に手すりをつけたい");
          formData.append("image", p.file);
        }
      }); 

      formData.append("comment", comment);

      const response = await fetch("/report", { method: "POST", body: formData });
      if (!response.ok) throw new Error("PDF生成に失敗しました");

      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      window.open(url);
      setTimeout(() => window.URL.revokeObjectURL(url), 1000);
    } catch (err) {
      console.error(err);
      alert("レポート生成に失敗しました");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="app-container">
      <PhotoUpload label="玄関の写真" photos={entrancePhotos} setPhotos={setEntrancePhotos} />
      <PhotoUpload label="廊下の写真" photos={hallwayPhotos} setPhotos={setHallwayPhotos} />

      <div className="comment-box">
        <h2>総評</h2>
        <textarea value={comment} onChange={(e) => setComment(e.target.value)} placeholder="総評を入力してください" />
      </div>

      <button className="report-button" onClick={handleGenerateReport} disabled={loading}>
        {loading ? "生成中..." : "レポート生成"}
      </button>
    </div>
  );
}
