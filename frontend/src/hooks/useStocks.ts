import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from '@/lib/api';
import { DailyRanking, StockDisplayData, StockInfo, RankingHistory } from '@/types/stock';

export const useLatestStocks = () => {
    return useQuery<DailyRanking[], Error, StockDisplayData[]>({
        queryKey: ['ranking', 'latest'],
        queryFn: async (): Promise<DailyRanking[]> => {
            const response = await api.get('/api/stocks/latest');
            return response.data;
        },
        select: (data: DailyRanking[]): StockDisplayData[] => {
            return data.map((ranking: DailyRanking) => ({
                rank: ranking.Rank,
                ticker: ranking.Stock.Ticker,
                name: ranking.Stock.Name,
                changeRate: ranking.ChangeRate,
                changeAmount: ranking.ChangeAmount,
                price: ranking.Price,
                aiAnalysis: ranking.AiAnalysis,
                date: ranking.Date,
            }));
        },
    });
};

export const useDateStocks = (date: string) => {
    return useQuery<DailyRanking[], Error, StockDisplayData[]>({
        queryKey: ['ranking', date],
        queryFn: async (): Promise<DailyRanking[]> => {
            const response = await api.get(`/api/stocks/date?date=${date}`);
            return response.data;
        },
        select: (data: DailyRanking[]): StockDisplayData[] => {
            return data.map((ranking: DailyRanking) => ({
                rank: ranking.Rank,
                ticker: ranking.Stock.Ticker,
                name: ranking.Stock.Name,
                changeRate: ranking.ChangeRate,
                changeAmount: ranking.ChangeAmount,
                price: ranking.Price,
                aiAnalysis: ranking.AiAnalysis,
                date: ranking.Date,
            }));
        },
        enabled: !!date,
    });
};

export const useStockHistory = (ticker: string) => {
    return useQuery<DailyRanking[], Error, { stockInfo: StockInfo; history: RankingHistory[] }>({
        queryKey: ['stocks', ticker],
        queryFn: async (): Promise<DailyRanking[]> => {
            const response = await api.get(`/api/stocks/${ticker}`);
            return response.data;
        },
        select: (data: DailyRanking[]): { stockInfo: StockInfo; history: RankingHistory[] } => {
            if (data.length === 0) {
                return {
                    stockInfo: {
                        ticker: ticker,
                        name: '',
                        sector: '',
                        industry: '',
                    },
                    history: [],
                };
            }

            const firstRanking = data[0];
            const stockInfo: StockInfo = {
                ticker: firstRanking.Stock.Ticker,
                name: firstRanking.Stock.Name,
                sector: firstRanking.Stock.Sector,
                industry: firstRanking.Stock.Industry,
            };

            const history: RankingHistory[] = data.map((ranking) => ({
                date: ranking.Date,
                rank: ranking.Rank,
                changeRate: ranking.ChangeRate,
                aiAnalysis: ranking.AiAnalysis,
            }));

            return {
                stockInfo,
                history,
            };
        },
        enabled: !!ticker,
    });
};

export const useSyncStocks = () => {
    const queryClient = useQueryClient();
    
    return useMutation({
        mutationFn: async (): Promise<void> => {
            await api.post('/api/admin/sync');
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['ranking'] });
        },
    });
};