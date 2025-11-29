'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Bell, HelpCircle } from 'lucide-react';

export default function Header() {
  const pathname = usePathname();

  return (
    <header className="sticky top-0 z-10 flex h-16 w-full items-center justify-center border-b border-zinc-200/50 bg-background-light/80 backdrop-blur-sm dark:border-zinc-800/50 dark:bg-background-dark/80">
      <div className="flex w-full max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
        <div className="flex items-center gap-4 text-zinc-900 dark:text-white">
          <div className="text-primary size-7">
            <svg fill="none" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg">
              <g clipPath="url(#clip0_6_330)">
                <path
                  clipRule="evenodd"
                  d="M24 0.757355L47.2426 24L24 47.2426L0.757355 24L24 0.757355ZM21 35.7574V12.2426L9.24264 24L21 35.7574Z"
                  fill="currentColor"
                  fillRule="evenodd"
                />
              </g>
              <defs>
                <clipPath id="clip0_6_330">
                  <rect fill="white" height="48" width="48" />
                </clipPath>
              </defs>
            </svg>
          </div>
          <Link href="/" className="text-lg font-bold leading-tight tracking-[-0.015em]">
            AI株分析
          </Link>
        </div>
        <div className="hidden items-center gap-8 md:flex">
          <Link
            href="/dashboard"
            className={`text-sm font-medium hover:text-primary dark:hover:text-primary ${
              pathname === '/dashboard'
                ? 'text-primary dark:text-primary'
                : 'text-zinc-600 dark:text-zinc-400'
            }`}
          >
            ダッシュボード
          </Link>
          <Link
            href="#"
            className="text-sm font-medium text-zinc-600 hover:text-primary dark:text-zinc-400 dark:hover:text-primary"
          >
            銘柄検索
          </Link>
          <Link
            href="/"
            className={`text-sm font-medium hover:text-primary dark:hover:text-primary ${
              pathname === '/'
                ? 'text-primary dark:text-primary'
                : 'text-zinc-600 dark:text-zinc-400'
            }`}
          >
            ランキング
          </Link>
        </div>
        <div className="flex items-center gap-2">
          <button className="flex h-10 w-10 cursor-pointer items-center justify-center overflow-hidden rounded-lg bg-zinc-100 text-zinc-600 transition-colors hover:bg-zinc-200 dark:bg-zinc-800 dark:text-zinc-300 dark:hover:bg-zinc-700">
            <Bell className="h-5 w-5" />
          </button>
          <button className="flex h-10 w-10 cursor-pointer items-center justify-center overflow-hidden rounded-lg bg-zinc-100 text-zinc-600 transition-colors hover:bg-zinc-200 dark:bg-zinc-800 dark:text-zinc-300 dark:hover:bg-zinc-700">
            <HelpCircle className="h-5 w-5" />
          </button>
        </div>
      </div>
    </header>
  );
}


