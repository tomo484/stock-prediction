'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { ChevronLeft, ChevronRight, Calendar } from 'lucide-react';

interface DateNavigationProps {
  currentDate?: string;
  onDateChange?: (date: string) => void;
  variant?: 'home' | 'dashboard';
  dateParam?: string | null;
}

export default function DateNavigation({
  currentDate,
  onDateChange,
  variant = 'home',
  dateParam,
}: DateNavigationProps) {
  const router = useRouter();
  const [selectedDate, setSelectedDate] = useState<string>(currentDate || '');

  useEffect(() => {
    if (dateParam) {
      setSelectedDate((prevDate) => {
        // 現在の値と異なる場合のみ更新（無限ループを防ぐ）
        return prevDate !== dateParam ? dateParam : prevDate;
      });
    } else if (currentDate) {
      // currentDateが指定されている場合
      setSelectedDate((prevDate) => {
        return prevDate !== currentDate ? currentDate : prevDate;
      });
    } else {
      // selectedDateが未設定の場合のみデフォルト値を設定
      setSelectedDate((prevDate) => {
        if (prevDate) return prevDate; // 既に値が設定されている場合は変更しない
        
        // デフォルトは今日の日付（米国時間基準の簡易実装）
        const today = new Date();
        const usDate = new Date(today.toLocaleString('en-US', { timeZone: 'America/New_York' }));
        const dateStr = usDate.toISOString().split('T')[0];
        return dateStr;
      });
    }
  }, [dateParam, currentDate]);

  // 日付文字列（YYYY-MM-DD）に日数を加減算するヘルパー関数
  // タイムゾーンの影響を受けない実装
  const addDaysToDateString = (dateStr: string, days: number): string => {
    if (!dateStr) return dateStr;
    
    // YYYY-MM-DD形式をパース
    const [year, month, day] = dateStr.split('-').map(Number);
    
    // UTC基準でDateオブジェクトを作成（タイムゾーンの影響を避ける）
    const date = new Date(Date.UTC(year, month - 1, day));
    
    // 日数を加減算
    date.setUTCDate(date.getUTCDate() + days);
    
    // YYYY-MM-DD形式に戻す
    const newYear = date.getUTCFullYear();
    const newMonth = String(date.getUTCMonth() + 1).padStart(2, '0');
    const newDay = String(date.getUTCDate()).padStart(2, '0');
    
    return `${newYear}-${newMonth}-${newDay}`;
  };

  const formatDate = (dateStr: string): string => {
    if (!dateStr) return '';
    try {
      const date = new Date(dateStr + 'T00:00:00');
      if (isNaN(date.getTime())) return dateStr;
      return date.toLocaleDateString('ja-JP', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        weekday: 'short',
      });
    } catch {
      return dateStr;
    }
  };

  const changeDate = (days: number) => {
    if (!selectedDate) return;
    
    // 日付文字列を直接操作してタイムゾーンの問題を回避
    const newDateStr = addDaysToDateString(selectedDate, days);
    
    // onDateChangeが指定されている場合は、それを使用（親コンポーネントが状態を管理）
    if (onDateChange) {
      onDateChange(newDateStr);
      setSelectedDate(newDateStr);
    } else {
      // router.pushを使用する場合、URLパラメータが更新されるので
      // useEffectがdateParamの変更を検知してselectedDateを更新する
      // そのため、ここではsetSelectedDateを呼ばない
      router.push(`?date=${newDateStr}`);
    }
  };

  const goToToday = () => {
    const today = new Date();
    const usDate = new Date(today.toLocaleString('en-US', { timeZone: 'America/New_York' }));
    const dateStr = usDate.toISOString().split('T')[0];
    setSelectedDate(dateStr);

    if (onDateChange) {
      onDateChange(dateStr);
    } else {
      router.push(`?date=${dateStr}`);
    }
  };

  if (variant === 'home') {
    return (
      <div className="flex items-center justify-between gap-2 rounded-lg bg-[#192233]/70 border border-[#324467] p-1.5 self-start sm:self-center">
        <button
          onClick={() => changeDate(-1)}
          className="flex h-8 w-8 items-center justify-center rounded-md text-slate-400 transition-colors hover:bg-white/10 hover:text-white"
        >
          <ChevronLeft className="h-5 w-5" />
        </button>
        <span className="text-white text-sm font-semibold whitespace-nowrap">
          {selectedDate ? formatDate(selectedDate) : '日付を選択'}
        </span>
        <button
          onClick={() => changeDate(1)}
          className="flex h-8 w-8 items-center justify-center rounded-md text-slate-400 transition-colors hover:bg-white/10 hover:text-white"
        >
          <ChevronRight className="h-5 w-5" />
        </button>
      </div>
    );
  }

  return (
    <div className="flex flex-wrap justify-between items-center gap-4 p-3 mb-6 bg-surface-light dark:bg-surface-dark rounded-lg border border-gray-200 dark:border-gray-700">
      <div className="flex items-center gap-2">
        <button
          onClick={() => changeDate(-1)}
          className="p-2 rounded-DEFAULT hover:bg-gray-100 dark:hover:bg-white/10"
        >
          <ChevronLeft className="h-5 w-5 text-text-secondary-light dark:text-text-secondary-dark" />
        </button>
        <button
          onClick={() => changeDate(1)}
          className="p-2 rounded-DEFAULT hover:bg-gray-100 dark:hover:bg-white/10"
        >
          <ChevronRight className="h-5 w-5 text-text-secondary-light dark:text-text-secondary-dark" />
        </button>
        <button
          onClick={goToToday}
          className="px-4 py-2 text-sm font-medium border border-gray-300 dark:border-gray-600 rounded-DEFAULT hover:bg-gray-100 dark:hover:bg-white/10"
        >
          今日
        </button>
      </div>
      <div className="flex items-center gap-4">
        <h2 className="text-lg font-bold">
          {selectedDate ? formatDate(selectedDate) : '日付を選択'}
        </h2>
      </div>
      <div className="flex items-center gap-2">
        <button className="p-2 rounded-DEFAULT hover:bg-gray-100 dark:hover:bg-white/10">
          <Calendar className="h-5 w-5 text-text-secondary-light dark:text-text-secondary-dark" />
        </button>
      </div>
    </div>
  );
}

