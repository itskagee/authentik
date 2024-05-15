# Generated by Django 5.0.6 on 2024-05-15 11:17

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ("authentik_stages_consent", "0006_alter_userconsent_expires"),
    ]

    operations = [
        migrations.AlterField(
            model_name="userconsent",
            name="expires",
            field=models.DateTimeField(db_index=True, default=None, null=True),
        ),
        migrations.AlterField(
            model_name="userconsent",
            name="expiring",
            field=models.BooleanField(db_index=True, default=True),
        ),
    ]
