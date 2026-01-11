import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { AuthProvider } from "@/contexts/AuthContext";
import { SnackbarProvider } from "@/contexts/SnackbarContext";
import Header from "@/components/Header";
import { getAuthToken, getAuthUsername, getAuthUserId } from "@/lib/authHelpers";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Meme9",
  description: "Social media platform",
  openGraph: {
    title: "Meme9",
    description: "Social media platform",
    type: "website",
  },
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  // Read auth state on server to prevent hydration mismatch
  const token = await getAuthToken();
  const username = await getAuthUsername();
  const userId = await getAuthUserId();
  
  const initialAuth = {
    isAuthenticated: !!(token && username && userId),
    username: username || null,
    userId: userId || null,
  };

  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable}`}
      >
        <AuthProvider initialAuth={initialAuth}>
          <SnackbarProvider>
            <Header />
            {children}
          </SnackbarProvider>
        </AuthProvider>
      </body>
    </html>
  );
}
