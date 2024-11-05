<template>
	<select @change="switchLanguage">
		<option
			v-for="locale in supportedLocales"
			:key="`locale-${locale.label}`"
			:value="locale.locale"
			:selected="currentLocale === locale.locale"
		>
			{{ locale.label}}
		</option>
	</select>
</template>

<script setup lang="ts">
import { AppTranslation, LOCALES, type LocaleTypes } from "@/core/config/locale.config"
import { useI18n } from 'vue-i18n'
import { useRouter } from "vue-router"

const {  locale: currentLocale } = useI18n<{ t: (key: string) => string; locale: string }>();

    const supportedLocales = [
      {label: "English", locale: LOCALES.EN},
      {label: "Русский", locale: LOCALES.RU},
      {label: "Қазақ", locale: LOCALES.KK},

    ];

    const router = useRouter();

    const switchLanguage = async (event: Event) => {
      const newLocale = (event.target as HTMLSelectElement).value as LocaleTypes;

      await AppTranslation.switchLanguage(newLocale);

      try {
        await router.replace({ params: { locale: newLocale } });
      } catch (e) {
        console.error("Error switching language:", e);
        await router.push("/");
      }
    };
</script>
