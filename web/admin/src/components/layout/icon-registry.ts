import {
  AlignHorizontalJustifyEnd,
  BookImage,
  Boxes,
  BriefcaseBusiness,
  ClipboardEdit,
  Component as ComponentIcon,
  Database,
  LaptopMinimalCheck,
  LucideColumnsSettings,
  LucideDatabaseZap,
  LucideHeadset,
  LucideMartini,
  LucideMessageSquareCode,
  LucideTableConfig,
  LucideTableProperties,
  LucideTags,
  LucideUserCog2,
  LucideUserRound,
  LucideUsers,
  MenuSquare,
  Network,
  TableOfContents,
  WholeWord,
  type LucideIcon,
} from 'lucide-react'

export type MenuIconOption = {
  name: string
  icon: LucideIcon
}

export const MENU_ICON_OPTIONS: MenuIconOption[] = [
  { name: 'component', icon: ComponentIcon },
  { name: 'TableOfContents', icon: TableOfContents },
  { name: 'BriefcaseBusiness', icon: BriefcaseBusiness },
  { name: 'AlignHorizontalJustifyEnd', icon: AlignHorizontalJustifyEnd },
  { name: 'LucideUsers', icon: LucideUsers },
  { name: 'LucideUserRound', icon: LucideUserRound },
  { name: 'MenuSquare', icon: MenuSquare },
  { name: 'Network', icon: Network },
  { name: 'LucideTableProperties', icon: LucideTableProperties },
  { name: 'LucideUserCog2', icon: LucideUserCog2 },
  { name: 'ClipboardEdit', icon: ClipboardEdit },
  { name: 'LaptopMinimalCheck', icon: LaptopMinimalCheck },
  { name: 'LucideHeadset', icon: LucideHeadset },
  { name: 'Boxes', icon: Boxes },
  { name: 'WholeWord', icon: WholeWord },
  { name: 'Database', icon: Database },
  { name: 'BookImage', icon: BookImage },
  { name: 'LucideColumnsSettings', icon: LucideColumnsSettings },
  { name: 'LucideTags', icon: LucideTags },
  { name: 'LucideMartini', icon: LucideMartini },
  { name: 'LucideDatabaseZap', icon: LucideDatabaseZap },
  { name: 'LucideMessageSquareCode', icon: LucideMessageSquareCode },
  { name: 'LucideTableConfig', icon: LucideTableConfig },
]

const MENU_ICON_REGISTRY: Record<string, LucideIcon> = {
  component: ComponentIcon,
  Component: ComponentIcon,
  'align-horizontal-justify-end': AlignHorizontalJustifyEnd,
  AlignHorizontalJustifyEnd,
  'book-image': BookImage,
  BookImage,
  boxes: Boxes,
  Boxes,
  'briefcase-business': BriefcaseBusiness,
  BriefcaseBusiness,
  'clipboard-edit': ClipboardEdit,
  ClipboardEdit,
  database: Database,
  Database,
  'laptop-minimal-check': LaptopMinimalCheck,
  LaptopMinimalCheck,
  'columns-settings': LucideColumnsSettings,
  LucideColumnsSettings,
  'database-zap': LucideDatabaseZap,
  LucideDatabaseZap,
  headset: LucideHeadset,
  LucideHeadset,
  martini: LucideMartini,
  LucideMartini,
  'message-square-code': LucideMessageSquareCode,
  LucideMessageSquareCode,
  'table-config': LucideTableConfig,
  LucideTableConfig,
  'table-properties': LucideTableProperties,
  LucideTableProperties,
  tags: LucideTags,
  LucideTags,
  'user-cog-2': LucideUserCog2,
  LucideUserCog2,
  'user-round': LucideUserRound,
  LucideUserRound,
  users: LucideUsers,
  LucideUsers,
  'menu-square': MenuSquare,
  MenuSquare,
  network: Network,
  Network,
  'table-of-contents': TableOfContents,
  TableOfContents,
  'whole-word': WholeWord,
  WholeWord,
}

function pascalToKebab(value: string): string {
  return value
    .replace(/^Lucide/, '')
    .replace(/([A-Z]+)([A-Z][a-z])/g, '$1-$2')
    .replace(/([a-z0-9])([A-Z])/g, '$1-$2')
    .toLowerCase()
}

export function getMenuIcon(iconName?: string): LucideIcon | undefined {
  const normalized = iconName?.trim()
  if (!normalized) return undefined

  return MENU_ICON_REGISTRY[normalized] ?? MENU_ICON_REGISTRY[pascalToKebab(normalized)]
}
