

export interface DailyRanking {
    ID: number;
    StockID: number;
    Date: string;
    Rank: number;
    Category: string;
    ChangeAmount: number;
    ChangeRate: number;
    Price: number;
    NewsSummary: string;
    AiAnalysis: string;
    Stock: Stock;
}

export interface Stock {
    ID: number;
    Ticker: string;
    Name: string;
    Sector: string;
    Industry: string;
    Ranking: DailyRanking[];
}

// 画面表示用の型定義
export type StockDisplayData = {
    rank: number;
    ticker: string;
    name: string;
    changeRate: number;
    changeAmount: number;
    price: number;
    aiAnalysis: string;
    date: string;
};

export type StockInfo = {
    ticker: string;
    name: string;
    sector: string;
    industry: string;
};

export type RankingHistory = {
    date: string;
    rank: number;
    changeRate: number;
    aiAnalysis: string;
};