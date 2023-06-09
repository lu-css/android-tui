package generate

import (
	"fmt"
	"strings"

	"github.com/lu-css/android-tui/src/files"
	"github.com/lu-css/android-tui/src/translate-xml"
	"github.com/lu-css/android-tui/src/utils"
	"github.com/lu-css/android-tui/src/validations"

	"github.com/manifoldco/promptui"
)

func baseTemplateLayout(activityName string) (string, string) {
	cleanActivityName := "activity_" + strings.Trim(activityName, "\n")

	template := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
    <androidx.constraintlayout.widget.ConstraintLayout xmlns:android="http://schemas.android.com/apk/res/android"
    xmlns:app="http://schemas.android.com/apk/res-auto"
    xmlns:tools="http://schemas.android.com/tools"
    android:layout_width="match_parent"
    android:layout_height="match_parent"
    tools:context=".%v">

    </androidx.constraintlayout.widget.ConstraintLayout>`, cleanActivityName)

	return template, cleanActivityName
}

func baseJavaLayout(ActivityName string, layoutName string) string {
	manifest := files.GetManifest()
	className := utils.CapitalizeFirstChar(ActivityName)

	javapackage := manifest.Package

	template := fmt.Sprintf(`
    package %s;

    import androidx.appcompat.app.AppCompatActivity;

    import android.os.Bundle;

    public class %s extends AppCompatActivity {

        @Override
        protected void onCreate(Bundle savedInstanceState) {
            super.onCreate(savedInstanceState);
            setContentView(R.layout.%s);
        }
    } `, javapackage, className, layoutName)

	return template

}
func ChooseActivity() error {
	activitiesTypes := []string{
		"Empty",
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U000027A1 {{ . | cyan }}",
		Inactive: "  {{ . | cyan }} ",
		Selected: "Activity Model: {{ . | red | cyan }}",
	}

	prompt := promptui.Select{
		Label:     "Select a Model",
		Items:     activitiesTypes,
		Templates: templates,
		Size:      4,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return nil
	}

	switch strings.ToLower(activitiesTypes[i]) {
	case "empty":
		genEmptyActivity()
		break
	}

	return nil
}

func genEmptyActivity() {
	prompt := promptui.Prompt{
		Label:    "Activity Name",
		Validate: validations.NonBlankInput,
	}

	activityName, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	layout, layoutName := baseTemplateLayout(activityName)
	javaCode := baseJavaLayout(activityName, layoutName)

	fmt.Println(layout)
	fmt.Println(javaCode)

	updateManifest(activityName)
}

func updateManifest(activityName string) {
	manifest := files.GetManifest()

	activity := translate_xml.Activity{
		MetaData: translate_xml.ActivityMetaData{},
		Exported: false,
		Name:     "." + activityName,
		Filter:   translate_xml.IntentFilter{},
	}

	manifest.Application.Activities = append(manifest.Application.Activities, activity)

	manifestFile := translate_xml.ToManifestFile(manifest)

	fmt.Println(manifestFile.GetXmlFile())
}
