import React from "react";

export function DefaultButton({ className, ...props }: React.ComponentProps<"button">) {
  return (
    <button
      data-slot="button"
      className={`bg-blue-500 cursor-pointer active:scale-95 h-9 px-4 py-2 has-[>svg]:px-3 text-white shadow-xs hover:bg-blue-600 inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4 shrink-0 [&_svg]:shrink-0 outline-none ${className}`}
      {...props}
    />
  );
}
