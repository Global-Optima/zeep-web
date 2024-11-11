import { toast, type ToastContainerOptions, type ToastOptions } from 'vue3-toastify'

export const toastConfig: ToastContainerOptions = {
	autoClose: 2000,
	theme: 'colored',
	position: 'top-center',
	dangerouslyHTMLString: true,
	hideProgressBar: true,
	pauseOnHover: false,
	icon: false,
	bodyClassName: 'text-base sm:text-xl',
	toastClassName: 'w-fit rounded-2xl px-4',
	containerClassName: 'flex flex-col justify-center items-center',
	closeButton: false,
}

export const toastSuccess = (msg: string, options?: ToastOptions) =>
	toast(msg, { type: 'success', position: 'top-center', ...options })

export const toastInfo = (msg: string, options?: ToastOptions) =>
	toast(msg, { type: 'info', position: 'top-center', ...options })
