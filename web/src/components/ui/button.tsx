import * as React from "react"
import { Slot } from "@radix-ui/react-slot"
import { cva, type VariantProps } from "class-variance-authority"

import { cn } from "@/lib/utils"

const buttonVariants = cva(
  "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring  disabled:cursor-not-allowed [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0",
  {
    variants: {
      variant: {
        default:
          "cursor-pointer bg-transparent text-gray-600  hover:bg-theme-color/70 hover:text-white text-sm rounded-full h-9 px-6 py-2 rounded m-1",
        destructive:
          "bg-destructive text-destructive-foreground shadow-sm hover:bg-destructive/90 h-9 px-6 py-2",
        outline:
          "border border-input bg-white-300 text-sm text-white shadow-sm hover:bg-white hover:text-theme-color/85 h-9 px-6 py-2",
        secondary:
          "bg-secondary text-theme-color shadow-sm hover:bg-secondary/80 hover:text-theme-color h-9 px-4 py-2",
        ghost: "hover:bg-accent text-gray-600 hover:text-gray-500 h-9 px-4 py-2",
        link: "text-primary underline-offset-4 hover:underline h-9 px-6 py-2" ,
        surely: "bg-theme-color text-sm text-gray-100 shadow hover:bg-theme-color/70 h-9 px-5 py-2",
        common:" cursor-pointer bg-theme-color text-white  hover:bg-theme-color/70 text-sm rounded-full h-9 lg:px-6 md:px-8 px-3 py-2 rounded m-1",
        cancel:" cursor-pointer border border-input bg-white text-sm text-gray-500 hover:bg-white hover:border-theme-color/50 hover:text-theme-color/85 rounded h-9 px-6 py-2  m-1",
        searchbtn:"cursor-pointer bg-theme-color text-white  hover:bg-theme-color/70 text-sm rounded-full h-7 lg:px-3 md:px-8 px-3 py-2 rounded m-1"
      },
      size: {
        // default: "h-9 px-3 py-2",
        sm: "h-8 rounded-md px-3 text-xs",
        lg: "h-10 rounded-md px-8",
        icon: "h-9 w-9",
      },
    },
    defaultVariants: {
      variant: "default",
      // size: "default",
    },
  }
)

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  asChild?: boolean
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, asChild = false, ...props }, ref) => {
    const Comp = asChild ? Slot : "button"
    return (
      <Comp
        className={cn(buttonVariants({ variant, size, className }))}
        ref={ref}
        {...props}
      />
    )
  }
)
Button.displayName = "Button"

export { Button, buttonVariants }
