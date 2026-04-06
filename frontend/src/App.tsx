import React, { useState } from "react";
import { PhotoUpload } from "./components/PhotoUpload";
import type { Photo } from "./components/PhotoUpload";

import "./index.css";

export default function App() {
  const [entrancePhotos, setEntrancePhotos] = useState<Photo[]>([
    { file: null, preview: "", text: "" },
    { file: null, preview: "", text: "" },
  ]);
  const [hallwayPhotos, setHallwayPhotos] = useState<Photo[]>([
    { file: null, preview: "", text: "" },
    { file: null, preview: "", text: "" },
  ]);
  const [loading, setLoading] = useState(false);

  const handleGenerateReport = async () => {
    try {
      setLoading(true);
      const formData = new FormData();

      entrancePhotos.forEach((p) => {
        if (p.file) {
          formData.append("entrance_texts", p.text);
          formData.append("entrance_images", p.file);
        }
      });

      hallwayPhotos.forEach((p) => {
        if (p.file) {
          formData.append("hallway_texts", p.text);
          formData.append("hallway_images", p.file);
        }
      });

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

      <button className="report-button" onClick={handleGenerateReport} disabled={loading}>
        {loading ? "生成中..." : "レポート生成"}
      </button>
    </div>
  );
}
