<template>
	<select @change="switchLanguage">
		<option
			v-for="sLocale in supportedLocales"
			:key="`locale-${sLocale}`"
			:value="sLocale"
			:selected="locale === sLocale"
		>
			{{ t(`locale.${sLocale}`) }}
		</option>
	</select>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { AppTranslation, type LocaleTypes } from "@/core/config/locale.config";
import { useI18n } from 'vue-i18n';
import { useRouter } from "vue-router";

export default defineComponent({
  setup() {
    // Types for locale
    const { t, locale } = useI18n<{ t: (key: string) => string; locale: string }>();

    // Infer type of supported locales from AppTranslation
    const supportedLocales: string[] = AppTranslation.supportedLocales;

    const router = useRouter();

    // Type for event in switchLanguage
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

    return { t, locale, supportedLocales, switchLanguage };
  }
});
</script>
