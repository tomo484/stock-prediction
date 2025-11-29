'use client';

import { Award, TrendingUp } from 'lucide-react';
import { RankingHistory } from '@/types/stock';

interface TimelineItemProps {
  history: RankingHistory;
  isFirst?: boolean;
  isLast?: boolean;
}

export default function TimelineItem({
  history,
  isFirst = false,
  isLast = false,
}: TimelineItemProps) {
  const formatDate = (dateStr: string): string => {
    try {
      const date = new Date(dateStr + 'T00:00:00');
      if (isNaN(date.getTime())) return dateStr;
      return date.toLocaleDateString('ja-JP', {
        year: 'numeric',
        month: 'numeric',
        day: 'numeric',
      });
    } catch {
      return dateStr;
    }
  };

  const Icon = isFirst ? Award : TrendingUp;
  const iconBgClass = isFirst ? 'bg-amber-400/20 text-amber-400' : 'bg-primary/20 text-primary';

  return (
    <>
      <div className="flex flex-col items-center gap-1.5 pt-1.5">
        <div className={`flex size-8 items-center justify-center rounded-full ${iconBgClass}`}>
          <Icon className="h-4 w-4" />
        </div>
        {!isLast && <div className="w-px grow bg-zinc-200 dark:bg-zinc-700" />}
      </div>
      <div className="flex flex-1 flex-col pb-10">
        <p className="text-xs text-zinc-500 dark:text-zinc-400">{formatDate(history.date)}</p>
        <div className="mt-2 flex flex-col gap-3 rounded-lg border border-zinc-200 bg-zinc-50 p-4 dark:border-zinc-800 dark:bg-zinc-800/50">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <p className="text-base font-bold text-zinc-900 dark:text-white">{history.rank}位</p>
            <p className="font-semibold text-green-500 dark:text-green-400">
              +{history.changeRate.toFixed(1)}%
            </p>
          </div>
          <div className="text-sm text-zinc-600 dark:text-zinc-300">
            <p className="font-medium">AI分析による急騰理由:</p>
            <p className="mt-1">{history.aiAnalysis}</p>
          </div>
        </div>
      </div>
    </>
  );
}

