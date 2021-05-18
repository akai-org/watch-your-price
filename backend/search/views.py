import requests
from django.http import HttpResponseBadRequest
from rest_framework import status
from rest_framework.response import Response
from rest_framework.views import APIView

from config.settings import SEARCH_MODULE_URL
from search.exceptions import NotEnoughParamsException


class SearchProducts(APIView):

    def get(self, request):
       
        try:
            params = self.__get_params(request.query_params)
        except NotEnoughParamsException as e:
            print(f'Failed to get params: {e}')
            return HttpResponseBadRequest()

        try:
            response = requests.get(SEARCH_MODULE_URL + '/search', params=params)
            response.raise_for_status()
        except requests.exceptions.RequestException:
            return HttpResponseBadRequest()

        return Response(response.json(), status=status.HTTP_200_OK)

    @staticmethod
    def __get_params(request_params):
        params = {}
        try:
            params['phrase'] = request_params['phrase']
            params['website'] = request_params['website']
        except KeyError:
            raise NotEnoughParamsException('Not enough params')
        params['page'] = request_params['page'] if 'page' in request_params else 0
        return params
