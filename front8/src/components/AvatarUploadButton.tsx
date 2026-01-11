"use client";

import { useState, useRef } from "react";
import { useRouter } from "next/navigation";
import { UsersClient } from "@/lib/api-clients";
import { SetAvatarRequest } from "@/schema/users";
import { getAuthToken } from "@/lib/authHelpers";
import styles from "./AvatarUploadButton.module.css";

interface AvatarUploadButtonProps {
  userId: string;
}

export default function AvatarUploadButton({
  userId,
}: AvatarUploadButtonProps) {
  const router = useRouter();
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileSelect = async (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    const file = event.target.files?.[0];
    if (!file) {
      return;
    }

    // Validate file type
    if (!file.type.startsWith("image/")) {
      setError("Please select an image file");
      return;
    }

    // Validate file size (e.g., max 5MB)
    const maxSize = 5 * 1024 * 1024; // 5MB
    if (file.size > maxSize) {
      setError("File size must be less than 5MB");
      return;
    }

    setUploading(true);
    setError(null);

    try {
      const token = await getAuthToken();
      const uploadResponse = await fetch(
        "https://meme2.mmaks.me/twirp/meme.photos.Photos/upload",
        {
          method: "POST",
          body: file,
          headers: {
            "Content-Type": file.type || "application/octet-stream",
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (!uploadResponse.ok) {
        const text = await uploadResponse.text().catch(() => "");
        throw new Error(
          `Failed to upload file: ${uploadResponse.status} ${
            uploadResponse.statusText
          }${text ? ` - ${text}` : ""}`
        );
      }

      const publicUrl = await uploadResponse.text();

      const setAvatarRequest = SetAvatarRequest.create({
        userId: userId,
        avatarUrl: publicUrl,
      });
      await UsersClient.SetAvatar(setAvatarRequest);

      // Success - refresh the page to show the new avatar
      router.refresh();
    } catch (err) {
      console.error("Avatar upload error:", err);
      setError(err instanceof Error ? err.message : "Failed to upload avatar");
    } finally {
      setUploading(false);
      if (fileInputRef.current) {
        fileInputRef.current.value = "";
      }
    }
  };

  const handleButtonClick = () => {
    fileInputRef.current?.click();
  };

  return (
    <div className={styles.container}>
      <input
        ref={fileInputRef}
        type="file"
        accept="image/*"
        onChange={handleFileSelect}
        className={styles.hiddenInput}
        disabled={uploading}
      />
      <button
        type="button"
        onClick={handleButtonClick}
        disabled={uploading}
        className={styles.button}
      >
        {uploading ? "Uploading..." : "Upload Avatar"}
      </button>
      {error && (
        <p className={styles.error}>{error}</p>
      )}
    </div>
  );
}
