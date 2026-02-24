import React from "react";

interface PhotoUploadProps {
  label: string;
  photos: Photo[];
  setPhotos: (photos: Photo[]) => void;
}

export interface Photo {
  file: File | null;
  preview: string;
}

export const PhotoUpload: React.FC<PhotoUploadProps> = ({ label, photos, setPhotos }) => {
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
    <div className="photo-upload">
      <h2>{label}</h2>
      <div style={{ display: "flex", flexDirection: "column", alignItems: "center", gap: "1rem" }}>
        {photos.map((p, i) => (
          <div key={i} className="photo-box">
            <input type="file" accept="image/*" onChange={(e) => handleFileChange(e, i)} />
            {p.preview ? <img src={p.preview} alt="preview" /> : <span>プレビューなし</span>}
          </div>
        ))}
      </div>
    </div>
  );
};
