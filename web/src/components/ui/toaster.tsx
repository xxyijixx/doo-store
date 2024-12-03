import { useToast } from "@/hooks/use-toast"
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import {
  Toast,
  ToastClose,
  ToastDescription,
  ToastProvider,
  ToastTitle,
  ToastViewport,
} from "@/components/ui/toast"
import successIcon from "@/assets/安装成功.png"
import falseIcon from "@/assets/安装失败.png"

export function SuccessToaster() {
  const { toasts } = useToast()

  return (
    <ToastProvider>
      {toasts.map(function ({ id, title, description, action, ...props }) {
        return (
          <Toast key={id} {...props}>
                <div className="grid gap-1">
                      <div className="flex items-center">
                        <Avatar className="size-5 mr-5">
                          <AvatarImage src={successIcon}/>
                          <AvatarFallback />
                        </Avatar>
                        {title && <ToastTitle className="font-bold">{title}</ToastTitle>}
                      </div>
                      {description && (
                        <ToastDescription className="ml-10 ">{description}</ToastDescription>
                      )}
                
                {action}
            </div>
            
            <ToastClose />
          </Toast>
        )
      })}
      <ToastViewport />
    </ToastProvider>
  )
}

export function FalseToaster() {
  const { toasts } = useToast()

  return (
    <ToastProvider>
      {toasts.map(function ({ id, title, description, action, ...props }) {
        return (
          <Toast key={id} {...props}>
                <div className="grid gap-1">
                  <div className="flex items-center">
                      <Avatar className="size-6 mr-4">
                        <AvatarImage src={falseIcon}/>
                        <AvatarFallback />
                      </Avatar>
                        {title && <ToastTitle className="font-bold">{title}</ToastTitle>}
                    </div>
                  {description && (
                    <ToastDescription className="ml-10">{description}</ToastDescription>
                  )}
                
                {action}
            </div>
            
            <ToastClose />
          </Toast>
        )
      })}
      <ToastViewport />
    </ToastProvider>
  )
}
