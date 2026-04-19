/// <reference types="vite/client" />

declare module 'is-hotkey' {
  function isHotkey(
    hotkey: string | string[],
    event?: { key: string; metaKey?: boolean; ctrlKey?: boolean; altKey?: boolean; shiftKey?: boolean }
  ): boolean
  export default isHotkey
}
