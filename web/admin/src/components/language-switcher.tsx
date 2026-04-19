import { Languages } from 'lucide-react'
import { useI18n } from '@/context/i18n-provider'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { useTranslation } from 'react-i18next'

export function LanguageSwitcher() {
  const { language, setLanguage } = useI18n()
  const { t } = useTranslation()

  const languages = [
    { value: 'zh', label: t('components.languageSwitcher.chinese') },
    { value: 'en', label: t('components.languageSwitcher.english') },
  ] as const

  const currentLanguage = languages.find((lang) => lang.value === language)

  return (
    <Select
      value={language}
      onValueChange={(value) => setLanguage(value as 'en' | 'zh')}
    >
      <SelectTrigger className='w-[140px]'>
        <Languages className='mr-2 h-4 w-4' />
        <SelectValue placeholder={currentLanguage?.label || t('components.languageSwitcher.chinese')} />
      </SelectTrigger>
      <SelectContent>
        {languages.map((lang) => (
          <SelectItem key={lang.value} value={lang.value}>
            {lang.label}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  )
}
