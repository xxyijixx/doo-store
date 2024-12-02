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
          "bg-theme-color text-white/80 shadow-sm hover:bg-secondary/80 hover:text-theme-color h-9 px-4 py-2",
        ghost: "hover:bg-accent text-gray-600 hover:text-gray-500 h-9 px-4 py-2",
        link: "text-primary underline-offset-4 hover:underline h-9 px-6 py-2",
        surely: "bg-theme-color text-sm text-gray-100 shadow hover:bg-theme-color/70 h-9 px-5 py-2",
        common: " cursor-pointer bg-theme-color text-white  hover:bg-theme-color/70 text-sm rounded-full h-9 lg:px-6 md:px-8 px-3 py-2 rounded m-1",
        cancel: " cursor-pointer border border-input bg-gray-200 text-sm text-gray-500 hover:bg-white hover:border-theme-color/50 hover:text-theme-color/85 rounded h-9 px-6 py-2  m-1",
        searchbtn: "cursor-pointer bg-transparent text-white  text-sm rounded h-7 lg:px-3 md:px-8 px-3 py-2 m-1",
        combar: " cursor-pointer  text-theme-color  border-theme-color/70  text-sm h-9 lg:px-0 md:px-0 px-0 py-2 m-1  rounded-none",
        defbar: "cursor-pointer bg-transparent text-gray-600  hover:text-theme-color/70  text-sm h-9 px-0 py-2  m-1",
        combarson: "cursor-pointer bg-theme-color/20 text-theme-color text-sm rounded-lg h-9 lg:px-5 md:px-8 px-2 py-2 m-1",
        defbarson: "cursor-pointer bg-transparent text-gray-600  hover:bg-theme-color/20 hover:text-theme-color text-sm rounded-lg h-9 px-5 py-2 m-1",
        insbtn:"border border-theme-color text-sm font-normal text-theme-color hover:text-theme-color/80 hover:border-theme-color/80 h-8 px-3 whitespace-nowrap cursor-pointer",
        indefbtn:"border border-input rounded-md bg-gray-300 text-sm text-white shadow-sm h-8 px-3 whitespace-nowrap cursor-not-allowed"
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
