import { ReactNode } from 'react';
import { Providers } from './providers';
import './globals.css';

export default function RootLayout({
  children,
}: {
  children: ReactNode;
}) {
  return (
    <html lang="ja" className="dark">
      <body className="font-display bg-background-light dark:bg-background-dark text-text-primary-light dark:text-text-primary-dark">
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}

