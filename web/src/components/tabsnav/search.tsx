// UniSearch.tsx
import { useState } from "react";
import { MagnifyingGlassIcon } from "@radix-ui/react-icons"
import { Input } from "@/components/ui/input";
import { Button } from '@/components/ui/button'
// import { debounce } from "lodash"; 
import { useTranslation } from "react-i18next";

interface UniSearchProps {
    onSearch: (query: string) => void; // 父组件传递的搜索函数

}

const UniSearch: React.FC<UniSearchProps> = ({ onSearch }) => {
    const { t } = useTranslation();
    const [query, setQuery] = useState(""); // 搜索框内容
    const [error, setError] = useState<string>(""); // 错误提示文字

    // 正则：只允许输入中文、英文、数字(包括')
    const regex = /^[a-zA-Z0-9\u4e00-\u9fa5']*$/;
    const chregex = /^[a-zA-Z0-9\u4e00-\u9fa5]*$/;

    // 处理搜索输入变化
    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const value = e.target.value;
        // 校验输入内容是否符合正则表达式
        if (regex.test(value)) {
            setQuery(value);
            setError(""); // 如果符合要求，清除错误提示
        } else {
            setError(t("请输入中文、英文或数字"));
        }
    };

    // 处理键盘输入的变化
    const handleKeyDown = (event: React.KeyboardEvent<HTMLInputElement>) => {
        // 只有在用户输入时才触发搜索（非中文输入时）
        if (event.key === "Enter") {
            handleSearch();
        }
    };

    // 处理提交搜索
    const handleSearch = () => {
        if (chregex.test(query) && !error) {
            onSearch(query); // 调用父组件传递的搜索函数
        } else {
            setError(t("请输入中文、英文或数字")); // 如果不符合 chregex，设置错误信息
        }
    };

    return (
        <div>
            <div className="relative flex group items-center w-10 h-10 hover:w-[200px] lg:mr-0 md:mr-0 -mr-9">

                <div 
                    className={`relative flex items-center h-10 bg-gray-200/50 rounded-full overflow-hidden transition-all duration-300 ${
                        query ? "w-[200px]" : "w-10 group-hover:w-[200px] focus-within:w-[200px]"
                    }`}
                >
                        <Input
                            type="text"
                            // placeholder={t("请输入中文、英文或数字...")}
                            value={query}
                            onChange={handleChange}
                            onKeyDown={handleKeyDown} // 键盘按下时触发（处理回车提交）
                            className="w-full h-full bg-transparent border-none pl-4 pr-10 focus:outline-none placeholder:text-gray-500"
                        >
                        </Input>
                        <Button
                                variant="searchbtn"
                                onClick={handleSearch}
                                disabled={!!error}
                                className="absolute lg:-right-0.5 md:-right-6 right-0 top-4 transform -translate-y-1/2 p-2 rounded-full"
                            >
                            <MagnifyingGlassIcon className="size-20 shrink-0 font-bold text-gray-800" />
                        </Button>
                </div>


            </div>
            {/* 错误提示文字 */}
            {error && (
                <div className="text-red-500 text-xs mt-1">{error}</div>
            )}
        </div>


    );
};

export default UniSearch;
