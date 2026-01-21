import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { AuthProvider } from "@/pkg/auth-context";
import { Toaster } from "sonner";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
    title: "BJDMS - Bangladesh Jatiotabadi Jubodal Management System",
    description: "Official management system for Bangladesh Jatiotabadi Jubodal",
    manifest: "/manifest.json",
};

export default function RootLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <html lang="en" className="dark">
            <body className={inter.className}>
                <AuthProvider>
                    <main className="min-h-screen bg-background">
                        {children}
                    </main>
                    <Toaster richColors position="top-right" />
                </AuthProvider>
            </body>
        </html>
    );
}
