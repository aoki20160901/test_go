import { useState } from "react";

interface Photo {
  file: File | null;
  preview: string;
}

export default function ReportForm() {
  const [entrancePhotos, setEntrancePhotos] = useState<Photo[]>([{ file: null, preview: "" }, { file: null, preview: "" }]);
  const [hallwayPhotos, setHallwayPhotos] = useState<Photo[]>([{ file: null, preview: "" }, { file: null, preview: "" }]);
  const [summary, setSummary] = useState("");
  const [pdfUrl, setPdfUrl] = useState<string | null>(null);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>, index: number, location: "entrance" | "hallway") => {
    if (!e.target.files) return;
    const file = e.target.files[0];
    const reader = new FileReader();
    reader.onloadend = () => {
      const newPhoto = { file, preview: reader.result as string };
      if (location === "entrance") {
        const updated = [...entrancePhotos];
        updated[index] = newPhoto;
        setEntrancePhotos(updated);
      } else {
        const updated = [...hallwayPhotos];
        updated[index] = newPhoto;
        setHallwayPhotos(updated);
      }
    };
    reader.readAsDataURL(file);
  };

  const handleSubmit = async () => {
    const formData = new FormData();
    entrancePhotos.forEach((p, i) => {
      if (p.file) {
        formData.append("text", "玄関に手すりをつけたい");
        formData.append("image", p.file);
      }
    });
    hallwayPhotos.forEach((p, i) => {
      if (p.file) {
        formData.append("text", "廊下に手すりをつけたい");
        formData.append("image", p.file);
      }
    });
    formData.append("text", summary);

    try {
      const response = await fetch("/report", {
        method: "POST",
        body: formData,
      });
      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      setPdfUrl(url);
    } catch (err) {
      console.error(err);
    }
  };

  const renderPhotoInputs = (photos: Photo[], location: "entrance" | "hallway") => (
    <div className="flex gap-4">
      {photos.map((p, i) => (
        <div key={i} className="flex flex-col items-center">
          <input type="file" accept="image/*" onChange={(e) => handleFileChange(e, i, location)} />
          {p.preview && <img src={p.preview} alt="preview" className="w-32 h-32 object-cover mt-2 border" />}
        </div>
      ))}
    </div>
  );

  return (
    <div className="bg-white p-6 rounded shadow-md w-full max-w-3xl">
      <h2 className="text-xl font-bold mb-2">玄関の写真</h2>
      {renderPhotoInputs(entrancePhotos, "entrance")}

      <h2 className="text-xl font-bold mt-4 mb-2">廊下の写真</h2>
      {renderPhotoInputs(hallwayPhotos, "hallway")}

      <h2 className="text-xl font-bold mt-4 mb-2">総評</h2>
      <textarea
        className="w-full border p-2 rounded"
        value={summary}
        onChange={(e) => setSummary(e.target.value)}
      />

      <button
        className="mt-4 bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
        onClick={handleSubmit}
      >
        レポート生成
      </button>

      {pdfUrl && (
        <div className="mt-6">
          <h2 className="text-lg font-bold mb-2">PDFプレビュー</h2>
          <iframe src={pdfUrl} className="w-full h-[600px] border" />
        </div>
      )}
    </div>
  );
}
