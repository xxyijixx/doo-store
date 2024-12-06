import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { 
    Avatar, 
    AvatarImage, 
    AvatarFallback 
} from '@/components/ui/avatar';
import { Item } from "@/type.d/common"
import { useState } from "react"
import * as http from "@/api/modules/fouceinter"
import { useTranslation } from "react-i18next"
import warnIcon from "@/assets/警告.png"
import { useToast } from "@/hooks/use-toast";
import { FalseToaster } from '../ui/toaster';


interface AlertDialogDemoProps {
    app: Item
    isOpen: boolean
    onClose: () => void
    onUninstall: () => void
}

export function AlertDialogDemo({ isOpen, onClose, app, onUninstall }: AlertDialogDemoProps) {
    const { t } = useTranslation();
    const { toast } = useToast();
    const [isLoading, setIsLoading] = useState(false);
    const [variantState, setVariantState] = useState<"success" | "destructive" | null>(null);

    // 卸载应用
    const handleUninstall = async () => {
        setIsLoading(true)
        try {
            // 发送 DELETE 请求来卸载应用
            const response = await http.deleteApp(app.key)
            if (response.code === 200) {
                onUninstall() // 卸载成功后，执行回调函数
                onClose() // 关闭对话框
            } else {
                alert(response.msg) // 异常情况，显示错误信息
            }
        } catch (error) {
            console.error(t("卸载应用失败:"), error);
            toast({
                title: t("卸载失败"),
                description: t("本次卸载失败，请重试。"),
                variant: "destructive",
                duration: 3000,
                className: "fixed top-16  lg:top-3 md:top-3 lg:right-6  md:right-4 right-1/2 translate-x-1/2 lg:translate-x-0 md:translate-x-0 w-[350px]"
            });
            setVariantState("destructive");
            
        } finally {
            setIsLoading(false)
        }
    }

    return (
        <>
        {variantState === "destructive" && <FalseToaster />}
        <AlertDialog open={isOpen} onOpenChange={(open) => { if (!open) onClose(); }}>
            <AlertDialogContent>
                <div className="flex items-center space-x-3">
                <Avatar>
                    <AvatarImage  src={warnIcon} />
                    <AvatarFallback>...</AvatarFallback>
                </Avatar>
                
                <AlertDialogHeader>
                    <AlertDialogTitle>{t('卸载')}</AlertDialogTitle>
                    <AlertDialogDescription>
                        {t('即将执行卸载操作，您是否确定要卸载此')} {app.name} {t('插件吗?')}
                    </AlertDialogDescription>
                </AlertDialogHeader>
                </div>
                
                    <AlertDialogFooter>
                            <AlertDialogCancel onClick={onClose}>{t('取消')}</AlertDialogCancel>
                            <AlertDialogAction
                                onClick={handleUninstall} // 在确认按钮上调用 handleUninstall
                                disabled={isLoading} // 如果正在加载，禁用按钮
                            >
                                {isLoading ? t("卸载中...") : t("确认")}
                            </AlertDialogAction>           
                    </AlertDialogFooter>
                
                
            </AlertDialogContent>
        </AlertDialog>
        </>
    )
}
