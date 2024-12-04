/* eslint-disable react-hooks/exhaustive-deps */
import { Button } from "@/components/ui/button"
import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
} from "@/components/ui/card"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { useState, useMemo } from "react"
import { AlertDialogDemo } from "@/components/tabsnav/uninstallalert"
import { LoadingOverlay } from "@/components/tabsnav/loading"
import {  AlertLogHave } from "@/components/tabsnav/logalert"
import EditDrawer from "@/components/drawer/editpage"
import { Item } from "@/type.d/common"
import { Skeleton } from "@/components/ui/skeleton"
import { useTranslation } from "react-i18next"
import { FalseToaster } from '@/components/ui/toaster'
import { useToast } from "@/hooks/use-toast";
import * as http from "@/api/modules/fouceinter";

interface InStalledBtnProps {
    app: Item;
    loadData: () => void;
}

export function InStalledBtn({ app, loadData }: InStalledBtnProps ) {

    const { t } = useTranslation()
    const { toast } = useToast();

    const [isDialogOpen, setIsDialogOpen] = useState(false)
    const [isLogHaveOpen, setIsLogHaveOpen] = useState(false)
    const [isDrawerOpen, setIsDrawerOpen] = useState(false)
    const openDialog = () => setIsDialogOpen(true)
    const openDrawer = () => {setIsDrawerOpen(true)}
    const closeDialog = () => setIsDialogOpen(false)
    const closeLogHave = () => setIsLogHaveOpen(false)
    const closeDrawer = () => setIsDrawerOpen(false)

    const [isLoading, setIsLoading] = useState(false)
    const handleUninstall = async () => {
        setIsLoading(true)
        await new Promise(resolve => setTimeout(resolve, 3000)) // 模拟卸载过程
        loadData(); // 卸载完成后触发 loadData 刷新数据
        setIsLoading(false)
        closeDialog() // 关闭对话框
    }
    const isRunning = useMemo(() => {
        return app.status == "Running"
    }, [app.status])

    const isInstalling = useMemo(() => {
        return app.status == "Installing"
    },[app.status])


    const appStatus = useMemo(() => {
        if(!app.status) {
            return t("未知")
        }
        switch(app.status) {
            case "Running":
                return t("已运行")
            case "Installing":
                return t("安装中")
            case "Stopped":
                return t("已停止")
            case "UpErr":
                return t("失败")
            case "Error":
                return t("错误")
            default:
                return t("未知")
        }
    }, [app.status])
    const handleToggleStarted = () => {
        let action = "stop"
        if (app.status == "Running") {
            action = "stop"
            setIsLoading(true)  // 开启全局 loading 蒙版
            http.putAppStatus(app.key, {
                action: action,
                params: {}
            }).then(res => {
                console.log(res)
                loadData()
                setIsLoading(false)  // 关闭 loading 蒙版
            }).catch(err => {
                console.error(err)
            })
        } else {
            if (app.status === "UpErr") {
                toast({
                    title: t("启动失败"),
                    description: t("请卸载后重新安装"),
                    variant: "destructive",
                    duration: 3000,
                    className: "fixed top-20 lg:top-3 md:top-3 lg:right-6  md:right-4 right-1/2 translate-x-1/2 lg:translate-x-0 md:translate-x-0 w-[350px]"
                });
                return;
            }
            action = "start"
            setIsLoading(true)  // 开启全局 loading 蒙版
            http.putAppStatus(app.key, {
                action: action,
                params: {}
            }).then(res => {
                console.log(res)
                    loadData()
                    setIsLoading(false)  // 关闭 loading 蒙版
            }).catch(err => {
                console.error(err)
                setIsLoading(false)
                toast({
                    title: t("错误提示"),
                    description: t("启动失败，请重试"),
                    variant: "destructive",
                    duration: 3000,
                    className: "fixed top-16 left-1/2 -translate-x-1/2 md:static md:translate-x-0 w-[90%] md:w-auto"
                })
            })
        }
    }

    const handleLogClick = () => {
        if (app.status != "Running") {
            // setIsLogDemoOpen(true)
            toast({
                title: t("温馨提示"),
                description: t("请先运行插件再查看日志~~~"),
            })
            
        } else {
            setIsLogHaveOpen(true)
        }
    }

    const statusDisplay = useMemo(() => {
        return (
            <span
                className={
                    !isRunning
                        ? "ml-3 border rounded-sm border-red-400 pt-1 px-2 text-sm font-normal text-red-400"
                        : "ml-3 border rounded-sm border-theme-color pt-1 px-2 text-sm font-normal text-theme-color"
                }
            >
                {appStatus}
            </span>
        );
    }, [isRunning, appStatus]);

    return (
        <>
        < FalseToaster/>
        <Card className="lg:w-auto  md:w-auto w-auto h-[180px] lg:mb-0 mb-3">
            <CardContent className="flex justify-start space-x-5 mt-6 pl-5">
                {isLoading ? (
                    <Skeleton className="h-12 w-12 rounded-full" />
                ) : (
                    <Avatar className="my-auto size-10">
                        <AvatarImage src={app.icon} />
                        <AvatarFallback>loading</AvatarFallback>
                    </Avatar>
                )}

                <CardDescription className="space-y-1 text-left">
                    {isLoading ? (
                        <Skeleton className="h-6 w-48" />
                    ) : (
                        <div className="text-xl font-medium text-slate-950 dark:text-white flex">
                            {app.name}
                            {statusDisplay}
                        </div>
                    )}

                    {isLoading ? (
                        <Skeleton className="h-4 w-56" />
                    ) : (
                        <p className="text-base line-clamp-2 min-h-[42px] leading-[21px] pt-1 pr-5">{app.description || t("No description available")}</p>
                    )}
                </CardDescription>
            </CardContent>
            <CardFooter className="flex justify-start -mt-1 gap-2 lg:gap-3 md:gap-3 lg:ml-14 md:ml-14 md:pl-6 pl-20">
                {isLoading ? (
                    <>
                        <Skeleton className="h-8 w-20 mx-2" />
                        <Skeleton className="h-8 w-20 mx-2" />
                        <Skeleton className="h-8 w-20 mx-2" />
                        <Skeleton className="h-8 w-20 mx-2" />
                    </>
                ) : (
                    <>
                        <Button
                            variant="insbtn"
                            className={`w-[56px] min-w-0 ${
                                !isRunning || isInstalling || isLoading
                                    ? "bg-gray-300 text-white cursor-not-allowed border-2 border-gray-300 hover:bg-gray-300 hover:text-white hover:border-2 hover:border-gray-300" 
                                    : ""
                            }`}
                            onClick={handleLogClick}
                            disabled={!isRunning || isInstalling || isLoading}
                        >
                            {t('日志')}
                        </Button>
                        {isLogHaveOpen && <AlertLogHave isOpen={isLogHaveOpen} onClose={closeLogHave} isLogOpen={false} app={app} />}

                        <Button 
                            variant="insbtn" 
                            onClick={openDrawer}  
                            disabled={!isRunning || isInstalling || isLoading}
                            className={`w-[56px] min-w-0 ${
                                !isRunning || isInstalling || isLoading ? "bg-gray-300 text-white cursor-not-allowed border-2 border-gray-300 hover:bg-gray-300 hover:text-white hover:border-2 hover:border-gray-300": ""
                            }`}
                        >
                            {t('参数')}
                        </Button>

                        <Button 
                            variant="insbtn" 
                            onClick={handleToggleStarted} 
                            disabled={isInstalling || isLoading}
                            className={`w-[56px] min-w-0 ${
                                !isRunning ? "border-theme-color text-theme-color" : ""
                            }
                            ${
                                isInstalling || isLoading ? "bg-gray-300 text-white cursor-not-allowed border-2 border-gray-300 hover:bg-gray-300 hover:text-white hover:border-2 hover:border-gray-300" : ""
                            }`}
                        >
                            {!isRunning ? t("启用") : t("停止")}
                        </Button>

                        <Button 
                            variant="insbtn"
                            onClick={openDialog}
                            disabled={isInstalling || isLoading}
                            className={`w-[56px] min-w-0 ${
                                isInstalling || isLoading ? "bg-gray-300 text-white cursor-not-allowed border-2 border-gray-300 hover:bg-gray-300 hover:text-white hover:border-2 hover:border-gray-300" : ""
                            }`}
                        >
                            {t('卸载')}
                        </Button>
                        <AlertDialogDemo 
                            isOpen={isDialogOpen} 
                            onClose={closeDialog} 
                            app={app} 
                            onUninstall={handleUninstall} 
                            
                            />

                        {/* <Button variant="common">重启</Button> */}
                    </>
                )}
            </CardFooter>

            {isLoading && <LoadingOverlay />}

            <EditDrawer isOpen={isDrawerOpen} onClose={closeDrawer} app={app} />
        </Card>
        </>
    )
}
