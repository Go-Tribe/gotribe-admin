import { Compass, ArrowLeft } from 'lucide-react'
import { useNavigate, useRouter } from '@tanstack/react-router'
import { Button } from '@/components/ui/button'
import { FeedbackState } from '@/components/feedback-state'

export function NotFoundError() {
  const navigate = useNavigate()
  const { history } = useRouter()
  return (
    <div className='flex min-h-svh items-center justify-center bg-background px-4 py-10'>
      <div className='w-full max-w-2xl'>
        <div className='mb-6 text-center'>
          <p className='text-xs font-semibold uppercase tracking-[0.2em] text-muted-foreground'>404</p>
          <h1 className='mt-2 text-4xl font-semibold tracking-tight sm:text-5xl'>Page not found</h1>
        </div>
        <FeedbackState
          title='This address is no longer available'
          description='The page may have been removed, renamed or never existed in this environment.'
          icon={<Compass className='h-5 w-5 text-primary' />}
          actions={
            <>
              <Button variant='outline' onClick={() => history.go(-1)}>
                <ArrowLeft className='mr-2 h-4 w-4' />
                Go Back
              </Button>
              <Button onClick={() => navigate({ to: '/' })}>Back to Home</Button>
            </>
          }
        />
      </div>
    </div>
  )
}
