import type { Metadata } from "next";
import Link from "next/link";
import "./globals.css";

export const metadata: Metadata = {
  title: "MuleLab — Outcome-accountable AI agents",
  description: "A commerce-agent sandbox where AI agents must prove measurable value.",
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en">
      <body>
        <header className="siteHeader">
          <Link className="brand" href="/"><span className="brandMark">M</span> MuleLab</Link>
          <nav aria-label="Primary navigation">
            <Link href="/demo">Run demo</Link>
            <Link href="/evals">Evals</Link>
            <Link href="/arena">Model arena</Link>
            <Link href="/architecture">Architecture</Link>
          </nav>
          <span className="mode"><i /> ZERO-COST MODE</span>
        </header>
        {children}
        <footer>Sticker Mule-inspired. Fictional data. No affiliation or private APIs.</footer>
      </body>
    </html>
  );
}

