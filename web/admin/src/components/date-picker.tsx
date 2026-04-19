import { format } from 'date-fns'
import { Calendar as CalendarIcon } from 'lucide-react'
import type { DateRange } from 'react-day-picker'
import { Button } from '@/components/ui/button'
import { Calendar } from '@/components/ui/calendar'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import { useTranslation } from 'react-i18next'

export type { DateRange }

type DatePickerProps = {
  selected: Date | undefined
  onSelect: (date: Date | undefined) => void
  placeholder?: string
  disabled?: boolean
}

export function DatePicker({
  selected,
  onSelect,
  placeholder,
  disabled = false,
}: DatePickerProps) {
  const { t } = useTranslation()
  const defaultPlaceholder = placeholder || t('components.datePicker.pickADate')
  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button
          variant='outline'
          data-empty={!selected}
          disabled={disabled}
          className='w-[240px] justify-start text-start font-normal data-[empty=true]:text-muted-foreground'
        >
          {selected ? (
            format(selected, 'MMM d, yyyy')
          ) : (
            <span>{defaultPlaceholder}</span>
          )}
          <CalendarIcon className='ms-auto h-4 w-4 opacity-50' />
        </Button>
      </PopoverTrigger>
      <PopoverContent className='w-auto p-0'>
        <Calendar
          mode='single'
          captionLayout='dropdown'
          selected={selected}
          onSelect={onSelect}
          disabled={(date: Date) =>
            date > new Date() || date < new Date('1900-01-01')
          }
        />
      </PopoverContent>
    </Popover>
  )
}

type DateRangePickerProps = {
  selected: DateRange | undefined
  onSelect: (range: DateRange | undefined) => void
  placeholder?: string
  disabled?: boolean
}

export function DateRangePicker({
  selected,
  onSelect,
  placeholder,
  disabled = false,
}: DateRangePickerProps) {
  const { t } = useTranslation()
  const defaultPlaceholder = placeholder || t('components.datePicker.pickARange')
  const label =
    selected?.from && selected?.to
      ? `${format(selected.from, 'yyyy-MM-dd')} ~ ${format(selected.to, 'yyyy-MM-dd')}`
      : selected?.from
        ? format(selected.from, 'yyyy-MM-dd')
        : null
  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button
          variant='outline'
          data-empty={!selected?.from}
          disabled={disabled}
          className='min-w-[260px] justify-start text-start font-normal data-[empty=true]:text-muted-foreground'
        >
          {label ?? <span>{defaultPlaceholder}</span>}
          <CalendarIcon className='ms-auto h-4 w-4 opacity-50' />
        </Button>
      </PopoverTrigger>
      <PopoverContent className='w-auto p-0' align='start'>
        <Calendar
          mode='range'
          captionLayout='dropdown'
          selected={selected}
          onSelect={onSelect}
          numberOfMonths={2}
          disabled={(date: Date) =>
            date > new Date() || date < new Date('1900-01-01')
          }
        />
      </PopoverContent>
    </Popover>
  )
}
