import { useEffect } from 'react'

export function usePageTitle(title: string) {
  useEffect(() => {
    const previousTitle = document.title

    if (title) {
      document.title = `${title} - Cloud Sentinel K8s`
    }

    return () => {
      document.title = previousTitle
    }
  }, [title])
}
