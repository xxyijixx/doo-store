import {
    Sheet,
    SheetContent,
    SheetHeader,
    SheetTitle,
    SheetDescription,
    SheetTrigger,
} from "@/components/ui/sheet"
import { ProfileForm } from "@/components/drawer/draform"
import  { useState } from "react";
import { Item } from "@/type.d/common";
import { useToast } from "@/hooks/use-toast";
import { Toaster } from '../ui/toaster';
import { useTranslation } from "react-i18next";

interface DrawerProps {
    status: string;
    isOpen: boolean;
    app: Item;
    loadData: () => void;
}

function Drawer({ status,app}: DrawerProps) {

    const { t } = useTranslation();
    const [isOpen, setIsOpen] = useState(false);
    const [buttonText, setButtonText] = useState(status === 'InUse' ? '已安装' : '安装'); // 用状态管理按钮文字
    const [currentStatus, setCurrentStatus] = useState(status); // 用状态管理当前的安装状态，便于控制按钮样式
    const { toast } = useToast();

    const handleInstallClick = () => {
        if (status === 'Unused') {
            setIsOpen(true); // 打开侧边栏
        }
        
    };

    const getButtonStyles = () => {
        if (currentStatus === 'InUse') {
            return 'border border-input rounded-md bg-gray-300 text-sm text-white shadow-sm  h-9 px-5 cursor-not-allowed '; // 在使用状态
        } else if (currentStatus === 'Unused') {
            return 'bg-theme-color text-sm text-gray-100 shadow hover:bg-theme-color/70 h-9 px-5'; // 未使用状态
        }
        return 'bg-theme-color text-white'; // 默认样式
    };


    //handleInstallSuccess成功运行安装局部更新状态
    const handleInstallSuccess = () => {
        // alert("安装成功");
        toast({
            title: "安装成功",
            description: "应用已成功安装",
            variant: "default",
        });
        setButtonText("已安装"); // 安装成功后更新按钮文字
        setCurrentStatus('InUse'); // 更新状态为 "InUse"（已安装），以改变按钮样式
        setIsOpen(false); // 关闭侧边栏
        
    };

    //handleInstallFalse失败运行安装局部更新状态
    const handleInstallFalse = () => {
        // alert("安装失败,请重试~");
        toast({
            title: "安装失败",
            description: "安装过程中发生错误，请重试。",
            variant: "destructive", // 失败提示通常使用 "destructive" 样式
        });
        setIsOpen(false); // 关闭侧边栏
    };
    
    

    return (
        <>
        <Toaster />
        <Sheet open={isOpen} onOpenChange={setIsOpen}>
            <SheetTrigger asChild>
                <div 
                onClick={(e) => {
                    if (status === 'InUse') {
                        e.preventDefault(); // 阻止默认行为
                        e.stopPropagation(); // 阻止事件冒泡
                    } else {
                        handleInstallClick(); // 处理安装逻辑
                    }
                }} 
                className={` rounded-md pt-2 ${getButtonStyles()}`}>
                    {buttonText}
                </div>
            </SheetTrigger>
            <SheetContent className='lg:overflow-y-hidden md:overflow-hidden overflow-auto '>
                <SheetHeader>
                    <SheetTitle className='ml-9 -mt-1.5 text-gray-700'>{t('返回')}</SheetTitle>
                    <SheetDescription className='pt-3'>
                    </SheetDescription>
                    <ProfileForm 
                        app={app} 
                        onInstallSuccess={handleInstallSuccess}
                        onFalse={handleInstallFalse}  
                        />
                </SheetHeader>
            </SheetContent>
        </Sheet>
        </>
    )
}



export default Drawer

