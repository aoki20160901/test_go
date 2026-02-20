import { useState } from "react";

interface Photo {
  file: File | null;
  preview: string;
}

function PhotoUpload({
  label,
  photos,
  setPhotos,
}: {
  label: string;
  photos: Photo[];
  setPhotos: (photos: Photo[]) => void;
}) {
  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>, i: number) => {
    if (!e.target.files) return;
    const file = e.target.files[0];
    const reader = new FileReader();
    reader.onloadend = () => {
      const updated = [...photos];
      updated[i] = { file, preview: reader.result as string };
      setPhotos(updated);
    };
    reader.readAsDataURL(file);
  };

  return (
    <div style={{ marginBottom: "1.5rem" }}>
      <h2>{label}</h2>
      <div style={{ display: "flex", flexDirection: "column", alignItems: "center", gap: "1rem" }}>
        {photos.map((p, i) => (
          <div
            key={i}
            style={{
              display: "flex",
              flexDirection: "column",
              alignItems: "center",
              gap: "0.5rem",
              padding: "0.75rem",
              border: "1px solid #ccc",
              borderRadius: "8px",
            }}
          >
            <input type="file" accept="image/*" onChange={(e) => handleFileChange(e, i)} />
            {p.preview ? (
              <img src={p.preview} alt="preview" style={{ width: 160, height: 160, objectFit: "cover" }} />
            ) : (
              <span>プレビューなし</span>
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
  const [hallPhotos, setHallPhotos] = useState<Photo[]>([
    { file: null, preview: "" },
    { file: null, preview: "" },
  ]);
  const [comment, setComment] = useState("");

  const handleGenerateReport = () => {
    alert("PDF生成（バックエンドと接続してください）");
  };

  return (
    <div style={{ padding: "1rem" }}>
      <PhotoUpload label="玄関の写真" photos={entrancePhotos} setPhotos={setEntrancePhotos} />
      <PhotoUpload label="廊下の写真" photos={hallPhotos} setPhotos={setHallPhotos} />

      <div style={{ marginBottom: "1rem" }}>
        <h2>総評</h2>
        <textarea
          value={comment}
          onChange={(e) => setComment(e.target.value)}
          rows={4}
          style={{ width: "100%", padding: "0.5rem", border: "1px solid #ccc", borderRadius: "4px" }}
        />
      </div>

      <button onClick={handleGenerateReport} style={{ padding: "0.75rem 1.5rem", cursor: "pointer" }}>
        レポート生成
      </button>
    </div>
  );
}
