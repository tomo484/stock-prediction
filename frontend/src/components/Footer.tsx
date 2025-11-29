export default function Footer() {
  return (
    <footer className="flex flex-col gap-6 px-5 py-10 text-center mt-16 border-t border-white/10">
      <div className="flex flex-wrap items-center justify-center gap-x-8 gap-y-4">
        <a
          href="#"
          className="text-slate-400 hover:text-white text-sm font-normal leading-normal transition-colors"
        >
          利用規約
        </a>
        <a
          href="#"
          className="text-slate-400 hover:text-white text-sm font-normal leading-normal transition-colors"
        >
          プライバシーポリシー
        </a>
        <a
          href="#"
          className="text-slate-400 hover:text-white text-sm font-normal leading-normal transition-colors"
        >
          お問い合わせ
        </a>
      </div>
      <p className="text-slate-500 text-xs font-normal leading-relaxed max-w-3xl mx-auto">
        © 2024 AI株分析. All rights reserved.
        <br />
        免責事項：本サービスが提供する情報は、投資判断の参考となる情報の提供を目的としたものであり、投資勧誘を目的としたものではありません。投資に関する最終決定は、ご自身の判断と責任において行ってください。
      </p>
    </footer>
  );
}


