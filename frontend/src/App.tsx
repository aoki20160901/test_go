// src/App.tsx
import React, { useState } from "react";

interface Photo {
  file: File | null;
  preview: string;
}

interface PhotoUploadProps {
  label: string;
  photos: Photo[];
  setPhotos: (photos: Photo[]) => void;
}

function PhotoUpload({ label, photos, setPhotos }: PhotoUploadProps) {
  const handleFileChange = (
    e: React.ChangeEvent<HTMLInputElement>,
    index: number
  ) => {
    if (!e.target.files) return;

    const file = e.target.files[0];
    const reader = new FileReader();

    reader.onloadend = () => {
      const updated = [...photos];
      updated[index] = {
        file,
        preview: reader.result as string,
      };
      setPhotos(updated);
    };

    reader.readAsDataURL(file);
  };

  return (
    <div style={{ marginBottom: "2rem", width: "100%" }}>
      <h2 style={{ fontSize: "1.2rem", fontWeight: "bold", marginBottom: "1rem" }}>
        {label}
      </h2>

      <div
        style={{
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          gap: "1rem",
        }}
      >
        {photos.map((p, i) => (
          <div
            key={i}
            style={{
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              gap: "0.5rem",
              padding: "1rem",
              border: "1px solid #ccc",
              borderRadius: "8px",
              width: "240px",
            }}
          >
            <input
              type="file"
              accept="image/*"
              onChange={(e) => handleFileChange(e, i)}
            />

            {p.preview ? (
              <img
                src={p.preview}
                alt="preview"
                style={{
                  width: 160,
                  height: 160,
                  objectFit: "cover",
                  borderRadius: "6px",
                }}
              />
            ) : (
              <span style={{ color: "#999" }}>プレビューなし</span>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}

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

    const response = await fetch("/report", {
      method: "POST",
      body: formData,
    });

    if (!response.ok) throw new Error("PDF生成に失敗しました");

    const blob = await response.blob();
    const url = window.URL.createObjectURL(blob);

    // Safari/iPhoneでも安全に表示
    const a = document.createElement("a");
    a.href = url;
    a.target = "_blank"; // 新しいタブで開く
    document.body.appendChild(a);
    a.click();
    a.remove();

    // メモリ解放
    setTimeout(() => window.URL.revokeObjectURL(url), 1000);
  } catch (error) {
    console.error(error);
    alert("レポート生成に失敗しました");
  } finally {
    setLoading(false);
  }
  };

  return (
    <div
      style={{
        maxWidth: "600px",
        margin: "0 auto",
        padding: "1.5rem",
        display: "flex",
        flexDirection: "column",
        gap: "1rem",
      }}
    >
      <PhotoUpload
        label="玄関の写真"
        photos={entrancePhotos}
        setPhotos={setEntrancePhotos}
      />

      <PhotoUpload
        label="廊下の写真"
        photos={hallwayPhotos}
        setPhotos={setHallwayPhotos}
      />

      <div>
        <h2 style={{ fontSize: "1.2rem", fontWeight: "bold" }}>総評</h2>
        <textarea
          value={comment}
          onChange={(e) => setComment(e.target.value)}
          rows={4}
          placeholder="総評を入力してください"
          style={{
            width: "100%",
            padding: "0.5rem",
            borderRadius: "6px",
            border: "1px solid #ccc",
            marginTop: "0.5rem",
          }}
        />
      </div>

      <button
        onClick={handleGenerateReport}
        disabled={loading}
        style={{
          padding: "0.75rem",
          backgroundColor: loading ? "#999" : "#2563eb",
          color: "white",
          border: "none",
          borderRadius: "6px",
          cursor: loading ? "not-allowed" : "pointer",
        }}
      >
        {loading ? "生成中..." : "レポート生成"}
      </button>
    </div>
  );
}
