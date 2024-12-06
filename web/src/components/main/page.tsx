/* eslint-disable @typescript-eslint/no-unused-vars */
/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable react-hooks/exhaustive-deps */
import { Skeleton } from "@/components/ui/skeleton"
import { Card, CardContent, CardDescription, CardFooter } from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import Drawer from "@/components/drawer/page";
import { useEffect, useState } from "react";
import { PaginationCom } from "@/components/tabsnav/paginationcom";
import UniSearch from "@/components/tabsnav/search"
import { InStalledBtn } from "@/components/tabsnav/installedbtn";
import { Button } from "@/components/ui/button";
// import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area"
import { Item } from "@/type.d/common";
import { Tag } from "@/api/interface/common"
import { motion, AnimatePresence } from "framer-motion"; // 引入 framer-motion
import * as http from '@/api/modules/fouceinter'
import { useTokenStore } from "@/store/ TokenContext";
import { useTranslation } from "react-i18next";
import { ChevronLeftIcon, ReloadIcon } from "@radix-ui/react-icons"

const POLLING_INTERVAL = 5000; // 5秒轮询一次

const fetchAppsData = async (tab: string, className = '', currentPage: number, pageSize = 9, query: string = '') => {

    const p = {
        page: 1,
        page_size: 9,
        name: '',
        description: '',
        class: ''
    }
    if (tab === 'all') {

        p.page = currentPage;
        p.page_size = pageSize;
        p.name = query;
        p.description = query;

        if (className && className !== 'all' && className !== 'allson') {
            p.class = className;
        }
    } else if (tab === 'installed') {
        p.page = currentPage;
        p.page_size = pageSize;
        p.name = query;
        p.description = query;

        if (className && className !== 'installed' && className !== 'allson') {
            p.class = className;
        }
    }

    try {
        let response;
        if (tab === 'all') {
            response = await http.getAppList(p)
        } else if (tab === 'installed') {
            response = await http.getInstalledAppList(p)
        }

        if (response?.code === 200) {
            if (response.data) {

                return {
                    items: response.data.items, // 返回应用数据
                    total: response.data.total, // 返回数据总数
                };

            }

        } else {
            throw new Error(response?.msg);
        }
    } catch (error) {
        console.error("Error fetching data:", error);
        return { items: [], total: 0 }; // 请求失败时返回空数组和总数为 0
    }
};


const handleMicroData = () => {
    const eventCenterForMicroApp = window.eventCenterForAppNameVite;
    //
    if (eventCenterForMicroApp) {
        // 设置名称
        // GlobalStore().setAppName(eventCenterForMicroApp.appName)
        // 主动获取基座下发的数据
        const info = eventCenterForMicroApp.getData();
        // if (info?.type == "init") {
        // initGlobaStore(info)
        // }
        useTokenStore.getState().token = info?.userInfo?.token
        console.log('22222222333333333', useTokenStore.getState().token);

        // 监听基座下发的数据变化
        eventCenterForMicroApp.addDataListener((data: any) => {
            // handleData(router, eventCenterForMicroApp.appName, data)
            console.log('1111111111', data);

            // useTokenStore().setToken(data.token)
        });
        // 监听基座下发的数据变化 - 全局
        eventCenterForMicroApp.addGlobalDataListener(() => {
            // useTokenStore().setToken(data.token)
        });
    }
}
handleMicroData();

function MainPage() {

    const { t } = useTranslation();

    const [apps, setApps] = useState<Item[]>([]);
    const [installedApps, setInstalledApps] = useState<Item[]>([]);  // 存储已安装的应用
    const [filteredApps, setFilteredApps] = useState<Item[]>([]);  // 存储过滤后的应用列表
    const [loading, setLoading] = useState(true);
    const [activeTab, setActiveTab] = useState("all"); // 当前激活的 Tab（"all" 或 "installed"）
    const [selectedClass, setSelectedClass] = useState("allson"); // 当前选中的 class 类型（如 "yundisk"）
    const [totalItems, setTotalItems] = useState(0); // 保存接口返回的总数
    const [currentPage, setCurrentPage] = useState(1); // 当前页面
    const [searchQuery, setSearchQuery] = useState(""); // 搜索关键词
    const [tags, setTags] = useState<Tag[]>([]);
    const pageSize = 9; // 每页显示的应用数
    const [isSearchExpanded, setIsSearchExpanded] = useState(false);


    // 请求数据
    const loadData = async (query: string = '', page: number = currentPage) => {
        setLoading(true);
        const data = await fetchAppsData(activeTab, selectedClass, page, pageSize, query);
        if (activeTab === 'installed') {
            console.log('installedApps', installedApps);
            setInstalledApps(data?.items || []); // 已安装应用的搜索结果
        } else if (activeTab === 'all') {
            console.log('apps', apps);
            setApps(data?.items || []); // 所有应用的搜索结果
        }
        setFilteredApps(data?.items || []); // 过滤后的应用
        setTotalItems(data?.total || 0); // 设置接口返回的总数
        setLoading(false);
    };

    

    const loadTags = async () => {
        const res = await http.getTags()
        if(res.data) {
            setTags(res.data)
        }
    }

    // 添加新的状态更新函数
    const updateAppsStatus = async () => {
        try {
            const data = await fetchAppsData(activeTab, selectedClass, currentPage, pageSize, searchQuery);
            if (data?.items) {
                // 只更新现有应用的状态
                setFilteredApps(prevApps => 
                    prevApps.map(prevApp => {
                        const updatedApp = data.items.find(newApp => newApp.id === prevApp.id);
                        return updatedApp ? { ...prevApp, status: updatedApp.status } : prevApp;
                    })
                );
            }
        } catch (error) {
            console.error('Error updating status:', error);
        }
    };

    // 修改轮询效果
    useEffect(() => {
        // 首次加载完整数据
        loadData(searchQuery, currentPage);

        // 只在 "installed" 标签页启用状态轮询
        if (activeTab === 'installed') {
            const intervalId = setInterval(() => {
                updateAppsStatus(); // 只更新状态
            }, POLLING_INTERVAL);

            return () => clearInterval(intervalId);
        }
    }, [activeTab, selectedClass, currentPage, searchQuery]);

    // 初次加载数据
    useEffect(() => {
        loadTags()
        loadData(searchQuery, currentPage); /// 根据 currentPage 加载数据



    }, [activeTab, selectedClass, currentPage, searchQuery]);


    // 切换 Tab 的方法
    const handleTabChange = (tab: string) => {
        setActiveTab(tab);
        setCurrentPage(1); // 切换 Tab 时重置为第 1 页
        setSearchQuery(""); // 如果切换到其他 Tab，清空搜索框

    };

    // 搜索时触发的过滤逻辑,父组件传递给 UniSearch 的搜索函数
    const handleSearch = (query: string) => {
        setSearchQuery(query); // 更新搜索关键词
        setCurrentPage(1); // 重置为第一页
        loadData(query); // 加载匹配查询的数据
    };

    // 处理分页变化
    const handlePageChange = (page: number) => {
        setCurrentPage(page); // 更新当前页
    };

    const totalPages = Math.ceil(totalItems / pageSize); // 确保分页计算是向上取整

    // 添加搜索框展开/收起的处理函数
    const handleSearchExpand = (expanded: boolean) => {
        setIsSearchExpanded(expanded);
    };

    return (
        <div className="flex flex-col min-h-[calc(100vh-66px)]">
            <div className="flex-none">
                <div className="flex justify-between items-center mb-4">
                    <div className="flex items-center justify-between flex-1">
                        <div className="flex items-center">
                            <Button 
                                variant="goback" 
                                size="icon"
                                onClick={() => window.history.back()}
                            >
                                <ChevronLeftIcon className="h-6 w-6" />
                            </Button>
                            
                            <h1 className="hidden sm:flex lg:text-3xl  text-2xl text-center text-gray-800">
                                {t('应用商店')}
                            </h1>

                            <Button 
                                variant="goback" 
                                size="icon"
                                onClick={() => loadData(searchQuery, currentPage)}
                                className="hidden md:block lg:block ml-2"
                            >
                                <ReloadIcon className="h-5 w-5" />
                            </Button>
                        </div>
                            <h1 className="sm:hidden text-2xl text-gray-800 lg:font-semibold md:text-3xl">
                                {t('应用商店')}
                            </h1>
                        <div className={`flex items-center justify-end ${isSearchExpanded ? 'flex-grow md:flex-grow-0 lg:flex-grow-0' : ''}`}>
                            <UniSearch 
                                onSearch={handleSearch} 
                                clearAfterSearch={false}
                                defaultValue={searchQuery}
                                onExpandChange={handleSearchExpand}
                            />
                        </div>
                    </div>
                </div>
                <AnimatePresence mode="wait">
                    <div key="b1" className="flex lg:-space-x-1 border-b lg:border-gray-200 md:border-gray-200 relative mb-3">
                        <>
                            <motion.div
                                key="Aoading"
                                initial={{ opacity: 0 }}
                                animate={{ opacity: 1 }}
                                exit={{ opacity: 0 }}
                                transition={{ duration: 1 }}
                                className="w-full lg:w-auto md:w-auto"
                            >
                                <ul className="flex items-center w-full">
                                    <li 
                                        className={`text-md pt-2 transition-all duration-300 relative z-10 w-1/2 lg:w-auto md:w-auto lg:mr-8 md:mr-8 ${
                                            activeTab === 'all' ? 'border-b-2 border-theme-color' : 'border-b-2 border-transparent'
                                        }`}
                                    >
                                        <Button
                                            variant={activeTab === "all" ? "combar" : "defbar"}
                                            onClick={() => handleTabChange("all")}
                                            className="w-full lg:w-auto md:w-auto lg:min-w-[10px] md:min-w-[10px] min-w-[130px] lg:justify-start md:justify-start justify-center"
                                        >
                                            {t('全部')}
                                        </Button>
                                    </li>
                                    <li
                                        className={`text-md pt-2 transition-all duration-300 relative z-10 w-1/2 lg:w-auto md:w-auto ${
                                            activeTab === 'installed' ? 'border-b-2 border-theme-color' : 'border-b-2 border-transparent'
                                        }`}
                                    >
                                        <Button
                                            variant={activeTab === "installed" ? "combar" : "defbar"}
                                            onClick={() => handleTabChange("installed")}
                                            className="w-full lg:w-auto md:w-auto lg:min-w-[10px] md:min-w-[10px] min-w-[130px] lg:justify-start md:justify-start justify-center"
                                        >
                                            {t('已安装')}
                                        </Button>
                                    </li>
                                </ul>
                            </motion.div>
                        </>
                    </div>
                </AnimatePresence>
            </div>

            <div className="flex-1 overflow-hidden">
                <div className="overflow-x-auto [&::-webkit-scrollbar]:hidden [-ms-overflow-style:'none'] [scrollbar-width:'none']">
                    <div 
                        className="flex space-x-2 mb-3 py-1 min-w-max cursor-grab active:cursor-grabbing"
                        onMouseDown={(e) => {
                            const ele = e.currentTarget;
                            const startX = e.pageX - ele.offsetLeft;
                            const scrollLeft = ele.parentElement?.scrollLeft || 0;
                            const handleMouseMove = (e: MouseEvent) => {
                                const x = e.pageX - ele.offsetLeft;
                                const walk = (x - startX) * 2;
                                if (ele.parentElement) {
                                    ele.parentElement.scrollLeft = scrollLeft - walk;
                                }
                            };
                            const handleMouseUp = () => {
                                document.removeEventListener('mousemove', handleMouseMove);
                                document.removeEventListener('mouseup', handleMouseUp);
                            };
                            document.addEventListener('mousemove', handleMouseMove);
                            document.addEventListener('mouseup', handleMouseUp);
                        }}
                    >
                        <Button
                            variant={selectedClass === "allson" ? "combarson" : "defbarson"}
                            onClick={() => setSelectedClass("allson")}
                        >
                            {t('全部')}
                        </Button>
                        {tags.map(tag => (
                            <Button 
                                key={tag.id}
                                variant={selectedClass === tag.key ? "combarson" : "defbarson"}
                                onClick={() => setSelectedClass(tag.key)}
                            >
                                {tag.name}
                            </Button>
                        ))}
                    </div>
                </div>

                <div className="mt-4">
                    <AnimatePresence mode="wait">
                        {/* 如果当前 Tab 是 "all" 或 "allson" 且未选择 class，显示 all 类应用列表 */}
                        {(activeTab === "all" && selectedClass !== "installed") && (
                            <div className={`grid lg:gap-4 md:gap-4 gap-2 m-1 grid-cols-1 md:grid-cols-2 lg:grid-cols-3 lg:max-h-[calc(100vh-280px)] md:max-h-[calc(100vh-280px)] max-h-[calc(100vh-230px)]  overflow-y-auto`}>
                                {loading ? (
                                    Array.from({ length: 9 }).map((_, index) => (
                                        <motion.div
                                            key={"a" + index}
                                            initial={{ opacity: 0 }}
                                            animate={{ opacity: 1 }}
                                            exit={{ opacity: 0 }}
                                            transition={{ duration: 0.5 }}
                                        >
                                            <Card key={index} className="lg:w-auto md:w-auto w-auto lg:h-[140px] lg:my-3 lg:mr-4 md:my-3 md:mr-4 my-3 mx-0 px-2">
                                                <CardContent className="flex justify-start space-x-4 pt-6">
                                                    <Skeleton className="h-10 w-10 rounded-full" />
                                                    <CardDescription className="space-y-2 text-left">
                                                        <Skeleton className="h-6 w-[200px] rounded-lg" />
                                                        <Skeleton className="h-4 w-[300px] rounded-lg" />
                                                    </CardDescription>
                                                </CardContent>
                                                <CardFooter className="flex justify-end pt-2">
                                                    <Skeleton className="h-8 w-[80px] rounded-lg" />
                                                </CardFooter>
                                            </Card>
                                        </motion.div>
                                    ))
                                ) : (
                                    filteredApps.map((app) => (
                                        <motion.div
                                            key={"d" + app.id}
                                            initial={{ opacity: 0 }}
                                            animate={{ opacity: 1 }}
                                            exit={{ opacity: 0 }}
                                            transition={{ duration: 0.7 }}
                                        >
                                            <Card key={app.id} className="lg:w-auto md:w-auto w-auto lg:h-[140px] lg:my-2 lg:mr-4 md:my-2 md:mr-4 my-0.5 mx-0 px-2">
                                                <CardContent className="flex flex-col mt-5 ">
                                                    <div className="flex w-full relative">
                                                        <div className="flex flex-1">
                                                            <Avatar className="my-auto mr-12 size-10">
                                                                <AvatarImage src={app.icon} />
                                                                <AvatarFallback>loading</AvatarFallback>
                                                            </Avatar>
                                                            <CardDescription className="space-y-1 text-left w-full">
                                                                <div className="lg:pr-0 md:pr-16 pr-16">
                                                                    <h1 className="text-xl font-medium text-slate-900 dark:text-white">{app.name}</h1>
                                                                </div>
                                                                <p className="text-base line-clamp-2 min-h-[42px] pt-1 w-11/12 md:w-4/5 lg:w-4/5">
                                                                    {app.description || "No description available"}
                                                                </p>
                                                            </CardDescription>
                                                        </div>
                                                        <div className="absolute right-0 -top-2">
                                                            <CardFooter className="pt-1">
                                                                <Drawer
                                                                    status={app.status}
                                                                    isOpen={false}
                                                                    app={app}
                                                                    loadData={loadData}
                                                                />
                                                            </CardFooter>
                                                        </div>
                                                    </div>
                                                </CardContent>
                                            </Card>
                                        </motion.div>
                                    ))
                                )}
                                <div className="mt-auto lg:hidden md:hidden">
                                    <PaginationCom
                                        currentPage={currentPage}
                                        totalPages={totalPages}
                                        totalItems={totalItems}
                                        pageSize={pageSize}
                                        onPageChange={handlePageChange}
                                        onPageSizeChange={(_: number) => {
                                            // 如果暂时不需要处理页面大小变化，可以留空
                                        }}
                                    />
                                </div>
                            </div>
                        )}


                        {/* 如果 Tab 是 "installed"，只显示已安装应用 */}
                        {activeTab === "installed" && (
                            <motion.div
                                key="eoading"
                                initial={{ opacity: 0 }}
                                animate={{ opacity: 1 }}
                                exit={{ opacity: 0 }}
                                transition={{ duration: 0.5 }}
                            >
                                <div className="grid lg:gap-6 md:gap-6 grid-cols-1 md:grid-cols-2 lg:grid-cols-3 lg:max-h-[calc(100vh-280px)] md:max-h-[calc(100vh-280px)] max-h-[calc(100vh-230px)] overflow-y-auto">
                                    {loading ? (
                                        <div></div>
                                    ) : (
                                        filteredApps.map((app) => (
                                            <motion.div
                                                key={"fo" + app.id}
                                                initial={{ opacity: 0 }}
                                                animate={{ opacity: 1 }}
                                                exit={{ opacity: 0 }}
                                                transition={{ duration: 0.7 }}
                                            >
                                                <InStalledBtn key={app.id} app={app} loadData={loadData} />
                                            </motion.div>
                                        ))
                                    )}
                                    <div className="mt-auto lg:hidden md:hidden">
                                        <PaginationCom
                                            currentPage={currentPage}
                                            totalPages={totalPages}
                                            totalItems={totalItems}
                                            pageSize={pageSize}
                                            onPageChange={handlePageChange}
                                            onPageSizeChange={(_: number) => {
                                                // 如果暂时不需要处理页面大小变化，可以留空
                                            }}
                                        />
                                    </div>
                                </div>
                            </motion.div>
                        )}
                    </AnimatePresence>
                </div>
            </div>

            <div className="flex-none">
                <div className="mt-auto hidden lg:block md:block">
                    <PaginationCom
                        currentPage={currentPage}
                        totalPages={totalPages}
                        totalItems={totalItems}
                        pageSize={pageSize}
                        onPageChange={handlePageChange}
                        onPageSizeChange={(_: number) => {
                            // 如果暂时不需要处理页面大小变化，可以留空
                        }}
                    />
                </div>
            </div>
        </div>
    )
}


export default MainPage; 