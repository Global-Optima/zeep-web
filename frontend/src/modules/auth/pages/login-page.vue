<script setup lang="ts">
import { ref, computed, type Ref } from 'vue'
import { Button } from '@/core/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/core/components/ui/card'
import { Input } from '@/core/components/ui/input'
import { Label } from '@/core/components/ui/label'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from '@/core/components/ui/select'
import { useRouter } from "vue-router"
import { getRouteName } from "@/core/config/routes.config"

interface Store {
  id: number;
  name: string;
}

interface Employee {
  id: number;
  name: string;
  storeId: number;
}

const stores: Ref<Store[]> = ref([
  { id: 1, name: 'Cafe A' },
  { id: 2, name: 'Cafe B' },
]);

const employees: Ref<Employee[]> = ref([
  { id: 1, name: 'Admin', storeId: 1 },
  { id: 2, name: 'Barista', storeId: 1 },
  { id: 3, name: 'Manager', storeId: 2 },
  { id: 4, name: 'Runner', storeId: 2 },
]);

const selectedStore = ref<string | undefined>(undefined);
const selectedEmployee = ref<string | undefined>(undefined);
const password = ref('');

const filteredEmployees = computed(() => {
  return employees.value.filter((employee) => employee.storeId === Number(selectedStore.value));
});

const isAdmin = ref(false);
const title = computed(() => (isAdmin.value ? 'Admin Login' : 'Employee Login'));
const description = computed(() => (isAdmin.value ? 'Enter your credentials to access the admin panel' : 'Enter your credentials to access the employee portal'));

const router = useRouter()

const login = () => {
  console.log('Logging in with:', {
    store: selectedStore.value,
    employee: selectedEmployee.value,
    password: password.value,
  });

  router.push({ name: getRouteName("ADMIN_DASHBOARD") })

};
</script>

<template>
	<div class="w-full h-full flex items-center justify-center ">
		<Card class="w-full max-w-sm">
			<CardHeader>
				<CardTitle class="text-2xl"> {{ title }} </CardTitle>
				<CardDescription> {{ description }} </CardDescription>
			</CardHeader>
			<CardContent>
				<div class="grid gap-6">
					<div class="grid gap-2">
						<Label for="store">Select Store</Label>
						<Select v-model="selectedStore">
							<SelectTrigger class="w-full">
								<SelectValue placeholder="Select a store" />
							</SelectTrigger>
							<SelectContent>
								<SelectGroup>
									<SelectLabel>Stores</SelectLabel>
									<SelectItem
										v-for="store in stores"
										:key="store.id"
										:value="store.id.toString()"
									>
										{{ store.name }}
									</SelectItem>
								</SelectGroup>
							</SelectContent>
						</Select>
					</div>

					<div
						v-if="selectedStore !== null"
						class="grid gap-2"
					>
						<Label for="employee">Select Employee</Label>
						<Select v-model="selectedEmployee">
							<SelectTrigger class="w-full">
								<SelectValue placeholder="Select an employee" />
							</SelectTrigger>
							<SelectContent>
								<SelectGroup>
									<SelectLabel>Employees</SelectLabel>
									<SelectItem
										v-for="employee in filteredEmployees"
										:key="employee.id"
										:value="employee.id.toString()"
									>
										{{ employee.name }}
									</SelectItem>
								</SelectGroup>
							</SelectContent>
						</Select>
					</div>

					<div class="grid gap-2">
						<Label for="password">Password</Label>
						<Input
							id="password"
							type="password"
							v-model="password"
							required
						/>
					</div>

					<Button
						type="button"
						class="w-full"
						@click="login"
					>
						Login
					</Button>
				</div>
			</CardContent>
		</Card>
	</div>
</template>
