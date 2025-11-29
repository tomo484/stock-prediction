'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { LayoutDashboard, Search } from 'lucide-react';

export default function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="w-64 flex-shrink-0 bg-surface-light dark:bg-surface-dark p-4 flex flex-col justify-between">
      <div className="flex flex-col gap-8">
        <div className="flex items-center gap-3 px-2">
          <div className="bg-center bg-no-repeat aspect-square bg-cover rounded-full size-10 bg-primary/20 flex items-center justify-center">
            <span className="text-primary font-bold text-lg">AI</span>
          </div>
          <h1 className="text-base font-bold">AI Stock Analyzer</h1>
        </div>
        <nav className="flex flex-col gap-2">
          <Link
            href="/dashboard"
            className={`flex items-center gap-3 px-3 py-2 rounded-DEFAULT transition-colors ${
              pathname === '/dashboard'
                ? 'bg-primary/20 text-primary'
                : 'text-text-secondary-light dark:text-text-secondary-dark hover:bg-gray-100 dark:hover:bg-white/10'
            }`}
          >
            <LayoutDashboard className="h-5 w-5" />
            <p className="text-sm font-medium">ダッシュボード</p>
          </Link>
          <Link
            href="#"
            className="flex items-center gap-3 px-3 py-2 rounded-DEFAULT text-text-secondary-light dark:text-text-secondary-dark hover:bg-gray-100 dark:hover:bg-white/10 transition-colors"
          >
            <Search className="h-5 w-5" />
            <p className="text-sm font-medium">銘柄検索</p>
          </Link>
        </nav>
      </div>
      <div className="flex items-center gap-3 p-2">
        <div className="bg-center bg-no-repeat aspect-square bg-cover rounded-full size-10 bg-primary/20 flex items-center justify-center">
          <span className="text-primary font-semibold text-sm">U</span>
        </div>
        <div className="flex flex-col">
          <p className="text-sm font-semibold">User</p>
          <p className="text-xs text-text-secondary-light dark:text-text-secondary-dark">設定</p>
        </div>
      </div>
    </aside>
  );
}


