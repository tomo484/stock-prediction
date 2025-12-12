import axios from 'axios';

const getBaseURL = () => {
  // 本番環境では環境変数から取得、なければ本番URLを使用
  if (process.env.NODE_ENV === 'production') {
    return process.env.NEXT_PUBLIC_API_URL || 'https://stock-prediction-zu8c.onrender.com';
  }
  // 開発環境ではローカルホストを使用
  return 'http://localhost:8080';
};

const api = axios.create({
  baseURL: getBaseURL(),
  headers: {
    'Content-Type': 'application/json',
  },
});

export default api;
