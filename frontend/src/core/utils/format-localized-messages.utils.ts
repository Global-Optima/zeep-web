export const formatLocalizedMessage = (s: string) =>
	s.replace(/\*(.*?)\*/g, '<span class="font-medium">"$1"</span>').replace(/\n/g, '<br>')
