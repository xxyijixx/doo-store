import { Button } from "@/components/ui/button"
import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
} from "@/components/ui/card"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { useState } from "react"
import { AlertDialogDemo } from "@/components/tabsnav/uninstallalert"
import { LoadingOverlay } from "@/components/tabsnav/loading"
import {  AlertLogHave } from "@/components/tabsnav/logalert"
import EditDrawer from "@/components/drawer/editpage"
import { Item } from "@/type.d/common"
import { Skeleton } from "@/components/ui/skeleton"
import { useTranslation } from "react-i18next"
import { FalseToaster } from '@/components/ui/toaster'
import { useToast } from "@/hooks/use-toast";

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

    const [isDisable, setIsDisable] = useState(false)
    const handleToggleStarted = () => {
        setIsDisable(!isDisable)
    }

    const handleLogClick = () => {
        if (isDisable) {
            // setIsLogDemoOpen(true)
            toast({
                title: t("温馨提示"),
                description: t("请先运行插件再查看日志~~~"),
            })
            
        } else {
            setIsLogHaveOpen(true)
        }
    }

    return (
        <>
        < FalseToaster/>
        <Card className="lg:w-auto  md:w-auto w-auto h-[180px] lg:mx-1 my-1 ">
            <CardContent className="flex justify-start space-x-5 mt-6">
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
                        <div className="text-xl font-semibold text-slate-950 dark:text-white flex">
                            {app.name}
                            <span
                                className={
                                    isDisable
                                        ? "ml-3  border rounded-sm border-red-400 pt-1 px-2 text-sm font-normal text-red-400"
                                        : "ml-3  border rounded-sm border-theme-color pt-1 px-2 text-sm font-normal text-theme-color"
                                }
                            >
                                {isDisable ? t("已停止") : t("已运行")}
                            </span>
                        </div>
                    )}

                    {isLoading ? (
                        <Skeleton className="h-4 w-56" />
                    ) : (
                        <p className="text-base line-clamp-2 min-h-[42px] leading-[21px] pt-1">{app.description || t("No description available")}</p>
                    )}
                </CardDescription>
            </CardContent>
            <CardFooter className="flex justify-start space-x-4  -mt-1 ml-14">
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
                            className={isDisable ? "bg-gray-300 text-white cursor-not-allowed border-2 border-gray-300 hover:bg-gray-300 hover:text-white hover:border-2 hover:border-gray-300" : ""}
                            onClick={handleLogClick}
                        >
                            {t('日志')}
                        </Button>
                        {isLogHaveOpen && <AlertLogHave isOpen={isLogHaveOpen} onClose={closeLogHave} isLogOpen={false} app={app} />}

                        <Button 
                            variant="insbtn" 
                            onClick={openDrawer}  
                            disabled={isDisable}  // 禁用按钮
                            className={isDisable ? "bg-gray-300 text-white cursor-not-allowed border-2 border-gray-300 hover:bg-gray-300 hover:text-white hover:border-2 hover:border-gray-300": ""}
                            
                            >
                            {t('参数')}
                        </Button>

                        <Button variant="insbtn" onClick={handleToggleStarted} className={isDisable ? "border-theme-color text-theme-color" : ""}>
                            {isDisable ? t("启用") : t("停止")}
                        </Button>

                        <Button variant="insbtn" onClick={openDialog}>
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
