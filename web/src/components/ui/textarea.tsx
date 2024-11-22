import * as React from "react"

import { cn } from "@/lib/utils"

export type TextareaProps = React.TextareaHTMLAttributes<HTMLTextAreaElement>

const Textarea = React.forwardRef<HTMLTextAreaElement, TextareaProps>(
  ({ className, ...props }, ref) => {
    return (
      <textarea
        className={cn(
          "flex lg:min-h-[200px] md:min-h-[200px] min-h-[150px] w-full rounded-md border border-input  hover:border-lime-400 bg-transparent px-3 py-2 text-sm shadow-sm placeholder:text-base focus-visible:outline-none focus-visible:ring-0 focus-visible:shadow-[0_0_10px_rgba(132,198,106,0.2),0_0_20px_rgba(132,198,106,0.2)] disabled:cursor-not-allowed disabled:opacity-50",
          className
        )}
        ref={ref}
        {...props}
      />
    )
  }
)
Textarea.displayName = "Textarea"

export { Textarea }
