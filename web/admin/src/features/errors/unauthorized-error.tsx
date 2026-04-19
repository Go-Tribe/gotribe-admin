import { ArrowLeft, ShieldAlert } from 'lucide-react'
import { useNavigate, useRouter } from '@tanstack/react-router'
import { Button } from '@/components/ui/button'
import { FeedbackState } from '@/components/feedback-state'

export function UnauthorisedError() {
  const navigate = useNavigate()
  const { history } = useRouter()
  return (
    <div className='flex min-h-svh items-center justify-center bg-background px-4 py-10'>
      <div className='w-full max-w-2xl'>
        <div className='mb-6 text-center'>
          <p className='text-xs font-semibold uppercase tracking-[0.2em] text-muted-foreground'>401</p>
          <h1 className='mt-2 text-4xl font-semibold tracking-tight sm:text-5xl'>Authorization required</h1>
        </div>
        <FeedbackState
          title='You need a valid session to continue'
          description='Please sign in with the right account, or go back if you reached this page by mistake.'
          tone='danger'
          icon={<ShieldAlert className='h-5 w-5 text-destructive' />}
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
