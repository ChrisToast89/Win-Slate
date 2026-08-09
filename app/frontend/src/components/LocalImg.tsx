// LocalImg — loads disk paths via the host FileAsDataURL bridge (WebView2-safe).

import React, { useEffect, useState } from 'react'
import { localMediaSrc } from '../lib/host'

export default function LocalImg({
  path,
  alt = '',
  title,
  className,
  style
}: {
  path: string
  alt?: string
  title?: string
  className?: string
  style?: React.CSSProperties
}): React.JSX.Element {
  const [src, setSrc] = useState('')

  useEffect(() => {
    let alive = true
    void localMediaSrc(path).then((u) => {
      if (alive) setSrc(u)
    })
    return () => {
      alive = false
    }
  }, [path])

  if (!src) return <span className={className} style={{ ...style, opacity: 0.3 }} />
  return <img className={className} style={style} src={src} alt={alt} title={title ?? path} />
}
