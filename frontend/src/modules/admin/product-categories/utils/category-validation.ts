import { toTypedSchema } from '@vee-validate/zod';
import * as z from 'zod';
import { MachineCategory } from './category-options';

export const createCategorySchema = toTypedSchema(
  z.object({
    name: z.string().min(1, 'Введите название категории'),
    description: z.string().optional(),
    machineCategory: z.nativeEnum(MachineCategory, {
      message: 'Выберите категорию машины',
    }),
  })
);
