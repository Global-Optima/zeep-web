import type { PaginationParams } from '@/core/utils/pagination.utils'

export interface AdditiveCategoriesFilterQuery extends PaginationParams {
  showAll?: boolean;
  productSizeId?: number;
  isMultipleSelect?: boolean;
  search?: string;
}

export interface AdditiveFilterQuery extends PaginationParams {
  search?: string;
  minPrice?: number;
  maxPrice?: number;
  categoryId?: number;
  productSizeId?: number;
}

export interface AdditiveDTO {
  id: number;
  name: string;
  description: string;
  basePrice: number;
  imageUrl: string;
  size: string;
  category: {
      id: number;
      name: string;
      isMultipleSelect: boolean;
  };
}

export interface AdditiveCategoryItemDTO {
  id: number;
  name: string;
  description: string;
  price: number;
  imageUrl: string;
  size: string;
  categoryId: number;
}

export interface AdditiveCategoryDTO {
  id: number;
  name: string;
  description: string;
  additives: AdditiveCategoryItemDTO[];
  isMultipleSelect: boolean;
}

export interface CreateAdditiveCategoryDTO {
  name: string;
  description?: string;
  isMultipleSelect: boolean;
}

export interface UpdateAdditiveCategoryDTO {
  name?: string;
  description?: string;
  isMultipleSelect?: boolean;
}

export interface UpdateAdditiveDTO {
  name?: string;
  description?: string;
  price?: number;
  imageUrl?: string;
  size?: string;
  additiveCategoryId?: number;
}

export interface AdditiveCategoryResponseDTO {
  id: number;
  name: string;
  description: string;
  isMultipleSelect: boolean;
}

export interface CreateAdditiveDTO {
  name: string;
  description: string;
  price: number;
  imageUrl?: string;
  size: string;
  additiveCategoryId: number;
}
