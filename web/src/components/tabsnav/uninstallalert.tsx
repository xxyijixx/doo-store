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
import { Item } from "@/type.d/common"
import { useState } from "react"
import * as http from "@/api/modules/fouceinter"
import { useTranslation } from "react-i18next"

interface AlertDialogDemoProps {
    app: Item
    isOpen: boolean
    onClose: () => void
    onUninstall: () => void
}

export function AlertDialogDemo({ isOpen, onClose, app, onUninstall }: AlertDialogDemoProps) {
    const { t } = useTranslation()
    const [isLoading, setIsLoading] = useState(false)

    // 卸载应用
    const handleUninstall = async () => {
        setIsLoading(true)
        try {
            // 发送 DELETE 请求来卸载应用
            // const response = await fetch(`http://192.168.31.214:8080/api/v1/apps/${app.key}`, {
            //     method: "DELETE",
            //     headers: {
            //         "token": "YIG8ANC8q2QxFV_Gf8qwkPdBj2EpsqGqlfc3qvSdg7ksVkZcokOUtQn43XGK0NK3BXUDsyebUlpKIFKXISMXA6nB0kpNgtZ2Vus-0ALbiLKPW74oqXtwUlA_aJyQP-hq", 
            //     },
            //     body: JSON.stringify(requestBody), 
            // }
            // )

            const response = await http.deleteApp(app.key)

            // const data = await response()

            if (response.code === 200) {
                onUninstall() // 卸载成功后，执行回调函数
                onClose() // 关闭对话框

            } else {
                alert(response.msg) // 异常情况，显示错误信息
            }
        } catch (error) {
            console.error("卸载应用失败:", error)
            alert("卸载失败")
        } finally {
            setIsLoading(false)
        }
    }

    return (
        <AlertDialog open={isOpen} onOpenChange={(open) => { if (!open) onClose(); }}>
            <AlertDialogContent>
                <AlertDialogHeader>
                    <AlertDialogTitle>{t('卸载')}</AlertDialogTitle>
                    <AlertDialogDescription>
                        {t('即将执行卸载操作，您是否确定要卸载此')} {app.name} {t('插件吗?')}
                    </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                    <AlertDialogAction
                        onClick={handleUninstall} // 在确认按钮上调用 handleUninstall
                        disabled={isLoading} // 如果正在加载，禁用按钮
                    >
                        {isLoading ? "卸载中..." : "确认"}
                    </AlertDialogAction>
                    <AlertDialogCancel onClick={onClose}>{t('取消')}</AlertDialogCancel>
                </AlertDialogFooter>
            </AlertDialogContent>
        </AlertDialog>
    )
}
