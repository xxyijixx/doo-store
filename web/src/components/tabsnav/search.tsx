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
    // 正则：只允许输入中文、英文、数字
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
            <div className="relative flex items-center lg:w-[300px] w-[300px]">

                <div className="w-screen">
                    <Input
                        type="text"
                        placeholder={t("请输入中文、英文或数字...")}
                        value={query}
                        onChange={handleChange}
                        onKeyDown={handleKeyDown} // 键盘按下时触发（处理回车提交）
                        className="input-class lg:w-[300px]  sm:w-[230px]"
                    />
                </div>


                <Button
                    variant="searchbtn"
                    onClick={handleSearch}
                    disabled={!!error}
                    className="absolute right-0 top-0.5  p-2 md:p-2"
                >
                    <MagnifyingGlassIcon className="size-10 shrink-0" />
                </Button>
            </div>
            {/* 错误提示文字 */}
            {error && (
                <div className="text-red-500 text-xs mt-1">{error}</div>
            )}
        </div>


    );
};

export default UniSearch;
