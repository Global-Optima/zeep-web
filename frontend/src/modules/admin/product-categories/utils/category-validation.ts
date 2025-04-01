import { toTypedSchema } from '@vee-validate/zod'
import * as z from 'zod'

export const createCategorySchema = toTypedSchema(
  z.object({
    name: z.string().min(1, 'Введите название категории'),
    description: z.string().optional(),
    machineCategory: z.enum(['TEA', 'COFFEE', 'ICE_CREAM'], {
      message: 'Выберите категорию машины',
    }),
  })
)
