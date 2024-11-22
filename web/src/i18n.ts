import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import HttpBackend from 'i18next-http-backend';
import LanguageDetector from 'i18next-browser-languagedetector';

i18n
    .use(HttpBackend) // 使用 HTTP 后端加载翻译文件
    .use(LanguageDetector) // 自动检测用户语言
    .use(initReactI18next) // 绑定到 react-i18next
    .init({
        fallbackLng: 'en', // 默认语言
        supportedLngs: ['en', 'ko', 'zh'], // 支持的语言
        debug: true, // 调试模式
        interpolation: {
            escapeValue: false, // React 已经自动转义
        },
        backend: {
            loadPath: '/locales/{{lng}}.json', // 翻译文件路径
        },
    });

export default i18n;
