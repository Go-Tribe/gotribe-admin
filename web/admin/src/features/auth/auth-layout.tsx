type AuthLayoutProps = {
  children: React.ReactNode
}

export function AuthLayout({ children }: AuthLayoutProps) {
  return (
    <div className='relative flex min-h-svh flex-col items-center justify-center overflow-hidden bg-background px-4 py-10'>
      <div className='pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_top_left,rgba(244,114,182,0.12),transparent_28%),radial-gradient(circle_at_bottom_right,rgba(14,165,233,0.1),transparent_32%)]' />
      <div className='pointer-events-none absolute inset-x-0 top-0 h-64 bg-[linear-gradient(180deg,rgba(255,255,255,0.6),transparent)] dark:bg-[linear-gradient(180deg,rgba(255,255,255,0.04),transparent)]' />
      <div className='relative mx-auto w-full max-w-[430px]'>{children}</div>
    </div>
  )
}
