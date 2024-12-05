/* eslint-disable react-hooks/exhaustive-deps */
// UniSearch.tsx
import { useState, useCallback } from "react";
import { MagnifyingGlassIcon } from "@radix-ui/react-icons"
import { Input } from "@/components/ui/input";
import { Button } from '@/components/ui/button'
// import { debounce } from "lodash"; 
import { useTranslation } from "react-i18next";

interface UniSearchProps {
    onSearch: (query: string) => void; // 父组件传递的搜索函数
    clearAfterSearch?: boolean; // 添加新属性控制搜索后是否清空
    defaultValue?: string; // 添加这一行
    onExpandChange?: (expanded: boolean) => void;  // 添加这一行
}

const UniSearch: React.FC<UniSearchProps> = ({ 
    onSearch, 
    clearAfterSearch = false,
    defaultValue = "",
    onExpandChange 
}) => {
    const { t } = useTranslation();
    const [query, setQuery] = useState(defaultValue);
    const [error, setError] = useState<string>("");
    const [isExpanded, setIsExpanded] = useState(false);

    // 正则：只允许输入中文、英文、数字(包括')
    const regex = /^[a-zA-Z0-9\u4e00-\u9fa5']*$/;
    const chregex = /^[a-zA-Z0-9\u4e00-\u9fa5]*$/;

    // 处理搜索输入变化
    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const value = e.target.value;
        if (regex.test(value)) {
            setQuery(value);
            setError("");
        } else {
            setError(t("请输入中文、英文或数字"));
        }
    };

    // 处理键盘输入
    const handleKeyDown = (event: React.KeyboardEvent<HTMLInputElement>) => {
        if (event.key === "Enter" && !error) {
            handleSearch();
        }
    };

    // 处理搜索
    const handleSearch = useCallback(() => {
        if (chregex.test(query) && !error) {
            onSearch(query);
            if (clearAfterSearch) {
                setQuery("");
            }
        } else {
            setError(t("请输入中文、英文或数字"));
        }
    }, [query, error, onSearch, clearAfterSearch, t, chregex]);

    // 添加鼠标滑过处理
    const handleMouseEnter = () => {
        setIsExpanded(true);
        onExpandChange?.(true);
    };

    const handleMouseLeave = () => {
        if (!document.activeElement?.classList.contains('search-input')) {
            setIsExpanded(false);
            onExpandChange?.(false);
        }
    };

    return (
        <div 
            onMouseEnter={handleMouseEnter}
            onMouseLeave={handleMouseLeave}
        >
            <div className={`relative flex group items-center ${
                query.length > 0 || isExpanded ? "w-[180px]" : "w-36"
            } h-36 lg:mr-0 md:mr-0 mr-0 transition-all duration-300`}>
                <div 
                    className={`relative flex items-center h-36 bg-gray-200/50 rounded-full overflow-hidden transition-all duration-300 ${
                        query.length > 0 || isExpanded ? "w-[180px]" : "w-36"
                    }`}
                    onClick={() => {
                        if (!isExpanded) {
                            setIsExpanded(true);
                            onExpandChange?.(true);
                        }
                    }}
                >
                    <Input
                        type="text"
                        value={query}
                        onChange={handleChange}
                        onKeyDown={handleKeyDown}
                        onFocus={() => {
                            setIsExpanded(true);
                            onExpandChange?.(true);
                        }}
                        onBlur={() => {
                            if (query.length === 0) {
                                setIsExpanded(false);
                                onExpandChange?.(false);
                            }
                            if (chregex.test(query) && !error) {
                                handleSearch();
                            }
                        }}
                        className="w-full h-full bg-transparent border-none pl-4 pr-10 focus:outline-none placeholder:text-gray-500 search-input"
                    />
                    <Button
                        variant="searchbtn"
                        onClick={(e) => {
                            e.stopPropagation();
                            if (query.length > 0) {
                                handleSearch();
                            } else {
                                setIsExpanded(true);
                                onExpandChange?.(true);
                            }
                        }}
                        disabled={!!error}
                        className="absolute lg:-right-3 md:-right-8 right-[-7px] top-4 transform -translate-y-1/2 p-2 rounded-full"
                    >
                        <MagnifyingGlassIcon className="size-6 shrink-0 font-bold text-gray-800" />
                    </Button>
                </div>
            </div>
            {error && (
                <div className="text-red-500 text-xs mt-1">{error}</div>
            )}
        </div>
    );
};

export default UniSearch;
