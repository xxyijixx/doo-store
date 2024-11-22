import * as React from "react"

import { cn } from "@/lib/utils"

export type InputProps = React.InputHTMLAttributes<HTMLInputElement>

const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className, type, ...props }, ref) => {
    return (
      
          <input
            type={type}
            className={cn(
              "flex w-full rounded-sm border border-input bg-transparent hover:border-lime-400 px-3 py-1 mt-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-base focus-visible:outline-none focus-visible:ring-0 focus-visible:shadow-[0_0_10px_rgba(132,198,106,0.2),0_0_20px_rgba(132,198,106,0.2)]  disabled:cursor-not-allowed disabled:opacity-50",
              className
            )}
            ref={ref}
            {...props}
          />
    
      
    )
  }
)
Input.displayName = "Input"

export { Input }
