import { useEffect } from 'react'

export function UnsavedChangesGuard({
  hasUnsavedChanges,
}: {
  hasUnsavedChanges: boolean
}) {
  return (
    <>
      <NativeBeforeUnloadGuard hasUnsavedChanges={hasUnsavedChanges} />

      {/* <RouterNavigationGuard hasUnsavedChanges={hasUnsavedChanges} /> */}
    </>
  )
}

type NativeBeforeUnloadGuardProps = {
  hasUnsavedChanges: boolean
}

function NativeBeforeUnloadGuard({
  hasUnsavedChanges,
}: NativeBeforeUnloadGuardProps) {
  useEffect(() => {
    if (!hasUnsavedChanges) {
      return
    }

    const handleBeforeUnload = (event: BeforeUnloadEvent) => {
      event.preventDefault()
      event.returnValue = ''
    }

    return () => {
      window.removeEventListener('beforeunload', handleBeforeUnload)
    }
  }, [hasUnsavedChanges])

  return null
}
